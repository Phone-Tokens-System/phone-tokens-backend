package repository

import (
	"context"
	"fmt"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type SmsRepository struct {
	db *gorm.DB
}

func NewSmsRepository(db *gorm.DB) *SmsRepository {
	return &SmsRepository{db: db}
}

func (r *SmsRepository) SaveSms(ctc context.Context, smsResponse model.SmsResponse) error {
	err := r.db.WithContext(ctc).Save(&smsResponse).Error
	return err
}

func (r *SmsRepository) GetAllSms(ctx context.Context) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse).Error
	return smsResponse, err
}

func (r *SmsRepository) GetSmsByServiceId(ctx context.Context, serviceId string) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse, "service_id = ?", serviceId).Error
	return smsResponse, err
}

func (r *SmsRepository) GetSmsByToken(ctx context.Context, token string) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	fmt.Println(token)
	//err := r.db.WithContext(ctx).Find(&smsResponse, "token = ?", token).Error
	err := r.db.Debug().
		WithContext(ctx).
		Find(&smsResponse, "token = ?", token).
		Error

	return smsResponse, err
}

func (r *SmsRepository) GetSmsByServiceName(ctx context.Context, serviceName string) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse, "service_name = ?", serviceName).Error
	return smsResponse, err
}
