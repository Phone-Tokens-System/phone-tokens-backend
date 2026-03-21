package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"
	"time"
)

type UserProfileHandler struct {
	userProfileService *users.UserProfileService
}

func NewUserProfileHandler(userProfileService *users.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService}
}

// GetFilters godoc
// @Summary Получить доступные фильтры профиля
// @Description Возвращает конфигурацию фильтров для UI
// @Tags user-profile
// @Produce json
// @Success 200 {object} dto.FilterResponse
// @Router /api/v1/user-profile/filters [get]
func (h *UserProfileHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	resp := h.userProfileService.GetFilters()

	writeJSON(w, http.StatusOK, resp)
}

// SaveUserProfile godoc
// @Summary Создать профиль пользователя
// @Description Создает или обновляет профиль текущего пользователя
// @Tags user-profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UserProfileRequest true "Данные профиля"
// @Success 200 {object} model.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user-profile [post]
func (h *UserProfileHandler) SaveUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId := user.UserID

	var req dto.UserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := parseTime(req.BirthDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userProfile := model.UserProfile{
		UserId:    userId,
		BirthDate: t,
		Gender:    req.Gender,
		Country:   req.Country,
		Region:    req.Region,
		City:      req.City,
		Education: req.Education,
	}
	userProfile.UserId = userId
	err = h.userProfileService.SaveUserProfile(r.Context(), userProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, http.StatusOK, userProfile)
}

// DeleteUserProfile godoc
// @Summary Удалить профиль пользователя
// @Description Удаляет профиль текущего пользователя
// @Tags user-profile
// @Produce json
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user-profile [delete]
func (h *UserProfileHandler) DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId := user.UserID
	err = h.userProfileService.DeleteUserProfile(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUserProfile godoc
// @Summary Обновить профиль пользователя
// @Description Обновляет профиль текущего пользователя
// @Tags user-profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UserProfileRequest true "Данные профиля"
// @Success 200 {object} model.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user-profile [put]
func (h *UserProfileHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId := user.UserID
	var req dto.UserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := parseTime(req.BirthDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userProfile := model.UserProfile{
		UserId:    userId,
		BirthDate: t,
		Gender:    req.Gender,
		Country:   req.Country,
		Region:    req.Region,
		City:      req.City,
		Education: req.Education,
	}
	userProfile.UserId = userId

	err = h.userProfileService.UpdateUserProfile(r.Context(), userId, &userProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, http.StatusOK, userProfile)
}

// GetUserProfileById godoc
// @Summary Получить профиль текущего пользователя
// @Description Возвращает профиль авторизованного пользователя
// @Tags user-profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfile
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/user-profile/me [get]
func (h *UserProfileHandler) GetUserProfileById(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId := user.UserID
	u, err := h.userProfileService.GetUserProfileByUserId(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDTO := dto.UserProfile{BirthDate: u.BirthDate, Gender: u.Gender,
		Country: u.Country, Region: u.Region, City: u.City, Age: u.Age}
	writeJSON(w, http.StatusOK, userDTO)
}

// GetUserProfileByToken godoc
// @Summary Получить профиль по токену
// @Security BearerAuth
// @Description Возвращает профиль пользователя по токену
// @Tags user-profile
// @Accept json
// @Produce json
// @Param request body dto.TokenRequest true "Токен"
// @Success 200 {object} dto.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/agents/tokens/user-profile [post]
func (h *UserProfileHandler) GetUserProfileByToken(w http.ResponseWriter, r *http.Request) {
	var req dto.TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.userProfileService.GetUserProfileByToken(r.Context(), req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userDTO := dto.ToUserProfile(u)
	writeJSON(w, http.StatusOK, userDTO)
}

// GetUserProfilesByAgentID godoc
// @Summary Получить пользователей агента
// @Description Возвращает список пользователей, чьи токены привязаны к агенту
// @Tags agents
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserProfile
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/agents/tokens/user-profile [get]
func (h *UserProfileHandler) GetUserProfilesByAgentID(w http.ResponseWriter, r *http.Request) {
	agentID, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userProfiles, err := h.userProfileService.GetUserProfilesForAgent(r.Context(), agentID.UserID)
	if err != nil {
		return
	}

	userDtos := make([]dto.UserProfile, 0)
	for _, u := range userProfiles {
		userDtos = append(userDtos, *dto.ToUserProfile(&u))
	}
	writeJSON(w, http.StatusOK, userDtos)
}

// GetUserProfilesFilteredByAgentID godoc
// @Summary Фильтрация пользователей агента
// @Description Возвращает пользователей агента с применением фильтров
// @Tags agents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.FilterRequest true "Фильтры пользователей" example({"filters":{"gender":"male","age_from":"20","age_to":"30"}})
// @Success 200 {array} dto.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/agents/tokens/user-profile/filtered [post]
func (h *UserProfileHandler) GetUserProfilesFilteredByAgentID(w http.ResponseWriter, r *http.Request) {
	agentID, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req dto.FilterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userProfiles, err := h.userProfileService.FilterUserProfilesByAgentId(r.Context(), req, agentID.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userDtos := make([]dto.UserProfile, 0)
	for _, u := range userProfiles {
		userDtos = append(userDtos, *dto.ToUserProfile(&u))
	}
	writeJSON(w, http.StatusOK, userDtos)
}

func GetUserFromContext(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(userContextKey).(*UserClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return claims, nil
}

func parseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
