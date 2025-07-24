package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// MakePayment godoc
// @Summary      Make payment
// @Description  Make a new payment
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        payment  body      payment.MakePaymentRequest  true  "Payment data"
// @Success      201  {object}  payment.Payment
// @Failure      400  {string}  string  "Invalid input"
// @Failure      500  {string}  string  "Failed to Make Payment"
// @Security     ApiKeyAuth
// @Router       /v1/payments [post]
func (h *Handler) MakePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req Payment

	// Json Decoder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		fmt.Println(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	payment_id, err := h.service.MakePayment(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to make payment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/v1/payments/%d", payment_id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment_id)
}
