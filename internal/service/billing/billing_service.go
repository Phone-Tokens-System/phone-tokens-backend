package billing

import (
	"context"
	"errors"
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"

	"gorm.io/gorm"
)

type BillingService struct {
	userRepo        users.Repository
	usageRepo       *repository.UsageRepository
	transactionRepo *repository.TransactionRepository
	token           string
	secret          string
}

func NewBillingService(u users.Repository, us *repository.UsageRepository, t *repository.TransactionRepository, token, secret string) *BillingService {
	return &BillingService{
		userRepo:        u,
		usageRepo:       us,
		transactionRepo: t,
		token:           token,
		secret:          secret,
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

// списание за SMS
func (s *BillingService) ChargeSms(ctx context.Context, userID string, cost float64) error {
	return s.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
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

		usage := &model.Usage{
			AgentID: userID,
			Service: "sms",
			Units:   1,
			Cost:    cost,
		}
		if err = s.usageRepo.SaveUsage(ctx, tx, usage); err != nil {
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
