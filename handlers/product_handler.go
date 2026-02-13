package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/produk
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {

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

func (h *ProductHandler) Handler(w http.ResponseWriter, r *http.Request) {
	response := models.APIResponse{
		Products: map[string]string{
			"DELETE /api/produk/:id": "delete product",
			"POST   /api/produk":     "create product",
			"PUT    /api/produk/:id": "update product",
			"GET    /api/produk":     "list all product",
			"GET    /api/produk/:id": "get product by id",
			"GET    /health":         "health check",
		},
		Categories: map[string]string{
			"DELETE /api/categories/:id": "delete category",
			"POST   /api/categories":     "create category",
			"PUT    /api/categories/:id": "update category",
			"GET    /api/categories":     "list all categories",
			"GET    /api/categories/:id": "get category",
		},
		Environment: "production",
		Message:     "API endpoints",
		Version:     "v1.0.0",
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")

	products, err := h.service.GetAll(name)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// HandleProductByID - GET/PUT/DELETE /api/produk/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {

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

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid product ID",
		})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid product ID",
		})
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid product ID",
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
		"message": "Product deleted successfully",
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
