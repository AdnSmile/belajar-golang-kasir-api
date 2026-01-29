package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandlerCategories - GET /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"message": "Method not allowed",
		})
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	categories, err := h.service.GetAll()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {

	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// HandleCategoryByID - GET/PUT/DELETE /api/produk/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"message": "Method not allowed",
		})
	}
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid category ID",
		})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid category ID",
		})
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid category ID",
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
}
