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

// ChargeSms — списание за SMS.
// Сначала пытаемся снять из пакетов агента, при их отсутствии — с баланса.
// Всё (списание единиц пакета / баланса + сохранение usage) атомарно в одной транзакции.
func (s *BillingService) ChargeSms(ctx context.Context, userID string, cost float64, unitsUsed int) error {
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		pkg, err := s.chargePackageForSmsTx(ctx, tx, userID, unitsUsed)
		if err == nil {
			// пакет списан — сохраняем usage с нулевой денежной стоимостью
			_ = pkg // подавляем предупреждение компилятора
			usage := &model.Usage{
				AgentID: userID,
				Service: model.SMS,
				Units:   unitsUsed,
				Cost:    0,
			}
			return s.usageRepo.SaveUsage(ctx, tx, usage)
		}

		if !errors.Is(err, ErrNoPackageUnitsLeft) {
			return err
		}

		// пакетов нет — списываем с денежного баланса
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
		usage := &model.Usage{
			AgentID: userID,
			Service: model.SMS,
			Units:   unitsUsed,
			Cost:    cost,
		}
		return s.usageRepo.SaveUsage(ctx, tx, usage)
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
		return s.transactionRepo.SaveTransaction(ctx, tx, txn)
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

// показать все доступные пакеты
func (s *BillingService) GetPackages(ctx context.Context) ([]model.Package, error) {
	return s.pkgRepo.GetPackages(ctx)
}

// AddAgentPkg — покупка пакета агентом.
// Атомарно: списание баланса + запись транзакции + создание agent_package.
// Если любой шаг провалится — вся операция откатывается.
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
		durationDays = 30
	}

	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// 1. Получаем агента с блокировкой строки
		agent, err := s.userRepo.GetAgentForUpdate(ctx, tx, agentId)
		if err != nil {
			return err
		}
		if agent.Balance < pkg.Price {
			return ErrNotEnoughBalance
		}

		// 2. Списываем баланс
		agent.Balance -= pkg.Price
		if _, err := s.userRepo.UpdateAgent(ctx, tx, agent); err != nil {
			return err
		}

		// 3. Сохраняем транзакцию
		txn := &model.Transaction{
			AgentID: agentId,
			Amount:  pkg.Price,
			Type:    model.Debit,
		}
		if err := s.transactionRepo.SaveTransaction(ctx, tx, txn); err != nil {
			return err
		}

		// 4. Создаём запись пакета агента
		pkgAgent := &model.AgentPackages{
			AgentId:     agentId,
			PackageId:   pkgId,
			ServiceType: pkg.Service,
			Status:      "ACTIVE",
			UnitsTotal:  pkg.Units,
			UnitsUsed:   0,
			ExpiresAt:   time.Now().AddDate(0, 0, durationDays),
		}
		return s.agentPkgRepo.AddAgentPackageTx(ctx, tx, pkgAgent)
	})
}

// UseAgentPkg — публичный метод для внешних вызовов (без транзакции).
func (s *BillingService) UseAgentPkg(ctx context.Context, agentPkgId int, unitsUsed int) (*model.AgentPackages, error) {
	pkg, err := s.agentPkgRepo.GetAgentPackageById(ctx, agentPkgId)
	if err != nil {
		return nil, err
	}
	pkg.UnitsUsed += int64(unitsUsed)
	pkg.UnitsTotal -= int64(unitsUsed)
	return s.agentPkgRepo.UpdateAgentPackage(ctx, agentPkgId, pkg)
}

// useAgentPkgTx — списание единиц пакета внутри существующей транзакции.
func (s *BillingService) useAgentPkgTx(ctx context.Context, tx *gorm.DB, agentPkgId int, unitsUsed int) (*model.AgentPackages, error) {
	pkg, err := s.agentPkgRepo.GetAgentPackageByIdTx(ctx, tx, agentPkgId)
	if err != nil {
		return nil, err
	}
	pkg.UnitsUsed += int64(unitsUsed)
	pkg.UnitsTotal -= int64(unitsUsed)
	return s.agentPkgRepo.UpdateAgentPackageTx(ctx, tx, agentPkgId, pkg)
}

// GetPackagesByAgentId — все пакеты агента.
func (s *BillingService) GetPackagesByAgentId(ctx context.Context, agentId string) ([]model.AgentPackages, error) {
	return s.agentPkgRepo.GetAgentPackagesByAgentId(ctx, agentId)
}

// GetTransactions — история транзакций агента.
func (s *BillingService) GetTransactions(ctx context.Context, agentID string) ([]model.Transaction, error) {
	return s.transactionRepo.GetTransactionsByAgentID(ctx, agentID)
}

// CreatePackage — создание тарифного пакета (администратор).
func (s *BillingService) CreatePackage(ctx context.Context, pkg *model.Package) error {
	if pkg.Id == (uuid.UUID{}) {
		pkg.Id = uuid.New()
	}
	if pkg.DurationDays <= 0 {
		pkg.DurationDays = 30
	}
	return s.pkgRepo.AddPackage(ctx, pkg)
}

// DeletePackage — удаление тарифного пакета (администратор).
func (s *BillingService) DeletePackage(ctx context.Context, pkgId string) error {
	pkg, err := s.pkgRepo.GetPackageByID(ctx, pkgId)
	if err != nil {
		return err
	}
	return s.pkgRepo.DeletePackage(ctx, pkg)
}

// ChargePackageForSms — публичный метод, используется снаружи (без tx).
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
			_ = s.agentPkgRepo.SetPackageStatus(ctx, pkg.Id, "EXPIRED")
			continue
		}
		if pkg.UnitsTotal >= int64(unitsUsed) {
			return s.UseAgentPkg(ctx, pkg.Id, unitsUsed)
		}
	}
	return nil, ErrNoPackageUnitsLeft
}

// chargePackageForSmsTx — то же, но внутри переданной транзакции.
func (s *BillingService) chargePackageForSmsTx(ctx context.Context, tx *gorm.DB, agentId string, unitsUsed int) (*model.AgentPackages, error) {
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
			_ = s.agentPkgRepo.SetPackageStatus(ctx, pkg.Id, "EXPIRED")
			continue
		}
		if pkg.UnitsTotal >= int64(unitsUsed) {
			return s.useAgentPkgTx(ctx, tx, pkg.Id, unitsUsed)
		}
	}
	return nil, ErrNoPackageUnitsLeft
}
