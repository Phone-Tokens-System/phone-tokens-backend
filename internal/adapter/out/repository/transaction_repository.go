package repository

import (
	"context"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) SaveTransaction(ctx context.Context, tx *gorm.DB, txn *model.Transaction) error {
	return tx.Create(txn).Error
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, ID string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, ID).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetTransactionsByAgentID(ctx context.Context, agentID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.WithContext(ctx).Find(&transactions, "agent_id = ?", agentID).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
