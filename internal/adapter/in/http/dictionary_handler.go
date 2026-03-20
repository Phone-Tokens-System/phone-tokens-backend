package http

import (
	"net/http"
	"phone-tokens/internal/adapter/dto"
)

type DictionaryHandler struct {
}

// GetCountries godoc
// @Summary Получить список стран
// @Description список стран
// @Tags dictionary
// @Accept json
// @Produce json
// @Success 200 {array} dto.CountryDTO
// @Failure 404 {object} map[string]string
// @Router /api/v1/dictionary/countries [get]
func (h *DictionaryHandler) GetCountries(w http.ResponseWriter, r *http.Request) {
	var res []dto.CountryDTO
	for _, c := range dto.GetCountries() {
		res = append(res, dto.CountryDTO{ID: c.ID, Name: c.Name})
	}

	writeJSON(w, http.StatusOK, res)
}

// GetRegions godoc
// @Summary Получить список регионов страны по ид
// @Description список регионов
// @Tags dictionary
// @Param country query string true "ID страны"
// @Accept json
// @Produce json
// @Success 200 {array} dto.RegionDTO
// @Failure 404 {object} map[string]string
// @Router /api/v1/dictionary/regions [get]
func (h *DictionaryHandler) GetRegions(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")

	var res []dto.RegionDTO

	for _, c := range dto.GetRegions(country) {
		res = append(res, dto.RegionDTO{ID: c.ID, Name: c.Name})
	}

	writeJSON(w, http.StatusOK, res)
}

// GetCities godoc
// @Summary Получить список городов по ид региона и ид страны
// @Description список городов
// @Tags dictionary
// @Param country query string true "ID страны"
// @Param region query string true "ID региона"
// @Accept json
// @Produce json
// @Success 200 {array} dto.RegionDTO
// @Failure 404 {object} map[string]string
// @Router /api/v1/dictionary/cities [get]
func (h *DictionaryHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	region := r.URL.Query().Get("region")

	cities := dto.GetCities(country, region)

	writeJSON(w, http.StatusOK, cities)
}
