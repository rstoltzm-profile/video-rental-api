package customer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCustomers godoc
// @Summary      List customers
// @Description  get all customers
// @Tags         customers
// @Produce      json
// @Success      200  {array}  customer.Customer
// @Security     ApiKeyAuth
// @Router       /v1/customers [get]
func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customers, err := h.service.GetCustomers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(customers)
}

// GetCustomerByID godoc
// @Summary      Get customer by ID
// @Description  Get a customer by their ID
// @Tags         customers
// @Produce      json
// @Param        id   path      int  true  "Customer ID"
// @Success      200  {object}  customer.Customer
// @Failure      400  {string}  string  "Invalid customer ID"
// @Failure      404  {string}  string  "Customer not found"
// @Security     ApiKeyAuth
// @Router       /v1/customers/{id} [get]
func (h *Handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	customer, err := h.service.GetCustomerByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(customer)
}

// CreateCustomer godoc
// @Summary      Create customer
// @Description  Create a new customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer  body      customer.CreateCustomerRequest  true  "Customer data"
// @Success      201  {object}  customer.Customer
// @Failure      400  {string}  string  "Invalid input"
// @Failure      500  {string}  string  "Failed to create customer"
// @Security     ApiKeyAuth
// @Router       /v1/customers [post]
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateCustomerRequest

	// Json Decoder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	customer, err := h.service.CreateCustomer(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/v1/customers/%d", customer.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomerByID godoc
// @Summary      Delete customer
// @Description  Delete a customer by their ID
// @Tags         customers
// @Param        id   path      int  true  "Customer ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {string}  string  "Invalid customer ID"
// @Failure      500  {string}  string  "Failed to delete customer"
// @Security     ApiKeyAuth
// @Router       /v1/customers/{id} [delete]
func (h *Handler) DeleteCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteCustomerByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
