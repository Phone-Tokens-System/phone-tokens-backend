package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"phone-tokens/internal/model"
	"phone-tokens/internal/service/tokens"
	"phone-tokens/internal/service/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

var _ users.Repository = (*Storage)(nil)
var _ tokens.Repository = (*Storage)(nil)

func (r *Storage) Save(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Storage) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Storage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Storage) CreateToken(ctx context.Context, token *model.UserToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Storage) GetTokenByID(ctx context.Context, id string) (*model.UserToken, error) {
	var token model.UserToken

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tokens.ErrNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (r *Storage) UpdateToken(ctx context.Context, token *model.UserToken) error {
	result := r.db.WithContext(ctx).
		Model(&model.UserToken{}).
		Where("id = ?", token.ID).
		Updates(map[string]interface{}{
			"name":        token.Name,
			"permissions": toPermissionStrings(token.Permissions),
			"status":      token.Status,
			"expires_at":  token.ExpiresAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tokens.ErrNotFound
	}
	return nil
}

func (r *Storage) DeleteToken(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tokens.ErrNotFound
	}
	return nil
}

func toPermissionStrings(perms []model.TokenPermission) []string {
	result := make([]string, 0, len(perms))
	for _, p := range perms {
		result = append(result, string(p))
	}
	return result
}

func (r *Storage) SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error) {
	log.Println(r)
	err := r.db.WithContext(ctx).Save(&request).Error
	fmt.Println(err)
	log.Println(err)
	if err != nil {
		return model.CsrRequest{}, err
	}
	return request, nil
}

func (r *Storage) GetCsrRequest(ctx context.Context, ID int) (*model.CsrRequest, error) {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	return &request, err
}

func (r *Storage) UpdateCsrStatus(ctx context.Context, ID int, status string) error {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	if err != nil {
		return err
	}

	request.Status = status
	return r.db.WithContext(ctx).Save(&request).Error
}

func (r *Storage) GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error) {
	requests := []model.CsrRequest{}
	err := r.db.WithContext(ctx).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Storage) SaveAgentInfo(ctx context.Context, info model.ExternalAgentInfo) error {
	err := r.db.WithContext(ctx).Save(&info).Error
	return err
}

func (r *Storage) GetAgentInfo(ctx context.Context, csrID int) (*model.ExternalAgentInfo, error) {
	var info model.ExternalAgentInfo
	err := r.db.WithContext(ctx).First(&info, "csr_id = ?", csrID).Error
	return &info, err
}

func (r *Storage) GetUserIdFromToken(ctx context.Context, token string) (string, error) {
	var user model.UserToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", users.ErrNotFound
		}
	}
	return user.UserID, nil
}

func (r *Storage) GetNumberFromUserId(ctx context.Context, userId string) (string, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", users.ErrNotFound
		}
	}
	return user.Phone, nil
}

func (r *Storage) GetTokenByToken(ctx context.Context, token string) (*model.UserToken, error) {
	var tokenObj model.UserToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&tokenObj).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
	}
	return &tokenObj, nil
}

func (r *Storage) SaveSms(ctc context.Context, smsResponse model.SmsResponse) error {
	err := r.db.WithContext(ctc).Save(&smsResponse).Error
	return err
}

func (r *Storage) GetAllSms(ctx context.Context) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse).Error
	return smsResponse, err
}

func (r *Storage) GetSmsByServiceId(ctx context.Context, serviceId uuid.UUID) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse, "service_id = ?", serviceId).Error
	return smsResponse, err
}

func (r *Storage) GetSmsByPhoneNumber(ctx context.Context, phoneNumber string) (model.SmsResponse, error) {
	smsResponse := model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse, "number = ?", phoneNumber).Error
	return smsResponse, err
}

func (r *Storage) GetSmsByServiceName(ctx context.Context, serviceName string) ([]model.SmsResponse, error) {
	smsResponse := []model.SmsResponse{}
	err := r.db.WithContext(ctx).Find(&smsResponse, "service_name = ?", serviceName).Error
	return smsResponse, err
}
