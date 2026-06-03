package billing

import (
	"context"
	"errors"
	"log/slog"
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingService struct {
	userRepo        users.Repository
	usageRepo       *repository.UsageRepository
	transactionRepo *repository.TransactionRepository
	pkgRepo         *repository.PackageRepository
	agentPkgRepo    *repository.AgentPackageRepository
	token           string
	secret          string
	frontendURL     string
}

func NewBillingService(
	u users.Repository,
	us *repository.UsageRepository,
	t *repository.TransactionRepository,
	pkgRepo *repository.PackageRepository,
	agentPkgRepo *repository.AgentPackageRepository,
	token, secret, frontendURL string,
) *BillingService {
	return &BillingService{
		userRepo:        u,
		usageRepo:       us,
		transactionRepo: t,
		pkgRepo:         pkgRepo,
		agentPkgRepo:    agentPkgRepo,
		token:           token,
		secret:          secret,
		frontendURL:     frontendURL,
	}
}

func (s *BillingService) TopUpBalance(ctx context.Context, agentID string, amount float64) error {
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		if amount <= 0 {
			return errors.New("amount must be greater than zero")
		}
		agent, err := s.userRepo.GetAgentForUpdate(ctx, tx, agentID)
		if err != nil {
			return err
		}
		agent.Balance += amount
		if _, err := s.userRepo.UpdateAgent(ctx, tx, agent); err != nil {
			return err
		}
		txn := &model.Transaction{
			AgentID: agentID,
			Amount:  amount,
			Type:    model.Credit,
		}

		if err := s.transactionRepo.SaveTransaction(ctx, tx, txn); err != nil {
			return err
		}

		return nil
	})
}

func (s *BillingService) TopDownBalance(ctx context.Context, agentID string, amount float64) error {
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		if amount <= 0 {
			return errors.New("amount must be greater than zero")
		}
		agent, err := s.userRepo.GetAgentForUpdate(ctx, tx, agentID)
		if err != nil {
			return err
		}
		slog.Info("agent balance", agentID, amount)
		if agent.Balance < amount {
			return ErrNotEnoughBalance
		}
		agent.Balance -= amount
		if _, err := s.userRepo.UpdateAgent(ctx, tx, agent); err != nil {
			return err
		}
		txn := &model.Transaction{
			AgentID: agentID,
			Amount:  amount,
			Type:    model.Debit,
		}

		if err := s.transactionRepo.SaveTransaction(ctx, tx, txn); err != nil {
			return err
		}

		return nil
	})
}

// списание за SMS
// сначала пытаемся снять из пакетов агента, далее с его баланса
func (s *BillingService) ChargeSms(ctx context.Context, userID string, cost float64, unitsUsed int) error {
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		_, err := s.ChargePackageForSms(ctx, userID, unitsUsed)
		if err == nil {
			usage := &model.Usage{
				AgentID: userID,
				Service: model.SMS,
				Units:   unitsUsed,
				Cost:    0,
			}
			if err = s.usageRepo.SaveUsage(ctx, tx, usage); err != nil {
				return err
			}
			return nil
		}

		if !errors.Is(err, ErrNoPackageUnitsLeft) {
			return err
		}

		agent, err := s.userRepo.GetAgentForUpdate(ctx, tx, userID)
		if err != nil {
			return err
		}

		if agent.Balance < cost {
			return errors.New("insufficient balance")
		}

		agent.Balance -= cost
		if _, err := s.userRepo.UpdateAgent(ctx, tx, agent); err != nil {
			return err
		}

		txn := &model.Transaction{
			AgentID: userID,
			Amount:  cost,
			Type:    model.Debit,
			Service: model.SMS,
		}

		if err := s.transactionRepo.SaveTransaction(ctx, tx, txn); err != nil {
			return err
		}

		// сохраняем usage при оплате с баланса
		usage := &model.Usage{
			AgentID: userID,
			Service: model.SMS,
			Units:   unitsUsed,
			Cost:    cost,
		}
		if err := s.usageRepo.SaveUsage(ctx, tx, usage); err != nil {
			return err
		}

		return nil
	})
}

// списание за звонок (по минутам)
func (s *BillingService) ChargeCall(ctx context.Context, userID string, costPerMin float64, duration int) error {
	totalCost := costPerMin * float64(duration)
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		agent, err := s.userRepo.GetAgentForUpdate(ctx, tx, userID)
		if err != nil {
			return err
		}

		if agent.Balance < totalCost {
			return errors.New("insufficient balance")
		}

		agent.Balance -= totalCost
		if _, err := s.userRepo.UpdateAgent(ctx, tx, agent); err != nil {
			return err
		}

		usage := &model.Usage{
			AgentID: userID,
			Service: "call",
			Units:   duration,
			Cost:    totalCost,
		}
		if err = s.usageRepo.SaveUsage(ctx, tx, usage); err != nil {
			return err
		}

		txn := &model.Transaction{
			AgentID: userID,
			Amount:  totalCost,
			Type:    model.Debit,
			Service: model.Call,
		}
		if err = s.transactionRepo.SaveTransaction(ctx, tx, txn); err != nil {
			return err
		}

		return nil
	})
}

func (s *BillingService) GetBalance(ctx context.Context, agentID string) (float64, error) {
	transactions, err := s.transactionRepo.GetTransactionsByAgentID(ctx, agentID)
	if err != nil {
		return 0, err
	}

	var balance float64
	for _, t := range transactions {
		switch t.Type {
		case model.Credit:
			balance += t.Amount
		case model.Debit:
			balance -= t.Amount
		}
	}

	return balance, nil
}

// показать все доступные пакеты (100 смс в месяц, 10 смс)
func (s *BillingService) GetPackages(ctx context.Context) ([]model.Package, error) {
	return s.pkgRepo.GetPackages(ctx)
}

// AddAgentPkg agent buys new package
func (s *BillingService) AddAgentPkg(ctx context.Context, pkgId string, agentId string) error {
	pkg, err := s.pkgRepo.GetPackageByID(ctx, pkgId)
	if err != nil {
		return err
	}
	if pkg.Id == (uuid.UUID{}) {
		return errors.New("package not found")
	}

	durationDays := pkg.DurationDays
	if durationDays <= 0 {
		durationDays = 30 // месяц по умолчанию
	}

	err = s.TopDownBalance(ctx, agentId, pkg.Price)
	if err != nil {
		return err
	}

	pkgAgent := &model.AgentPackages{
		AgentId:     agentId,
		PackageId:   pkgId,
		ServiceType: pkg.Service, // копируем тип сервиса из пакета
		Status:      "ACTIVE",
		UnitsTotal:  pkg.Units,
		UnitsUsed:   0,
		ExpiresAt:   time.Now().AddDate(0, 0, durationDays),
	}
	err = s.agentPkgRepo.AddAgentPackage(ctx, pkgAgent)
	if err != nil {
		return err
	}
	return nil
}

// UseAgentPkg agent uses units of package
func (s *BillingService) UseAgentPkg(ctx context.Context, agentPkgId int, unitsUsed int) (*model.AgentPackages, error) {
	pkg, err := s.agentPkgRepo.GetAgentPackageById(ctx, agentPkgId)
	if err != nil {
		return nil, err
	}
	pkg.UnitsUsed += int64(unitsUsed)
	pkg.UnitsTotal -= int64(unitsUsed)
	newPkg, err := s.agentPkgRepo.UpdateAgentPackage(ctx, agentPkgId, pkg)
	if err != nil {
		return nil, err
	}
	return newPkg, nil
}

// GetPackagesByAgentId Get all packages agent owns
func (s *BillingService) GetPackagesByAgentId(ctx context.Context, agentId string) ([]model.AgentPackages, error) {
	return s.agentPkgRepo.GetAgentPackagesByAgentId(ctx, agentId)
}

// GetTransactions возвращает историю транзакций агента (для раздела «Аккаунтинг»)
func (s *BillingService) GetTransactions(ctx context.Context, agentID string) ([]model.Transaction, error) {
	return s.transactionRepo.GetTransactionsByAgentID(ctx, agentID)
}

// CreatePackage создаёт новый тарифный пакет (вызывает администратор)
func (s *BillingService) CreatePackage(ctx context.Context, pkg *model.Package) error {
	if pkg.Id == (uuid.UUID{}) {
		pkg.Id = uuid.New()
	}
	if pkg.DurationDays <= 0 {
		pkg.DurationDays = 30
	}
	return s.pkgRepo.AddPackage(ctx, pkg)
}

// DeletePackage удаляет тарифный пакет (вызывает администратор)
func (s *BillingService) DeletePackage(ctx context.Context, pkgId string) error {
	pkg, err := s.pkgRepo.GetPackageByID(ctx, pkgId)
	if err != nil {
		return err
	}
	return s.pkgRepo.DeletePackage(ctx, pkg)
}

func (s *BillingService) ChargePackageForSms(ctx context.Context, agentId string, unitsUsed int) (*model.AgentPackages, error) {
	agentPkgs, err := s.GetPackagesByAgentId(ctx, agentId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for _, pkg := range agentPkgs {
		if pkg.ServiceType != model.SMS {
			continue
		}
		if pkg.Status != "ACTIVE" {
			continue
		}
		if pkg.ExpiresAt.Before(now) {
			// пакет истёк — помечаем EXPIRED, не прерываем цикл
			_ = s.agentPkgRepo.SetPackageStatus(ctx, pkg.Id, "EXPIRED")
			continue
		}
		if pkg.UnitsTotal >= int64(unitsUsed) {
			agentPkg, err := s.UseAgentPkg(ctx, pkg.Id, unitsUsed)
			if err != nil {
				return nil, err
			}
			return agentPkg, nil
		}
	}
	return nil, ErrNoPackageUnitsLeft
}
