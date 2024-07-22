package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/common"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/service"

	"go.uber.org/zap"
)

type Handlers struct {
	Service *service.Service
	Logger  *zap.Logger
}

func NewHandlers(service *service.Service, logger *zap.Logger) *Handlers {
	return &Handlers{
		Service: service,
		Logger:  logger,
	}
}

type OrderRequest struct {
	OrderType service.OrderType `json:"orderType"`
	StockID   int               `json:"stockID"`
	Quantity  int               `json:"quantity"`
}

func (h *Handlers) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Health check request received.")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy"}`))
}

func (h *Handlers) OrderHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Order request received.")

	txnToken := r.Header.Get("Txn-Token")

	ctx := context.WithValue(r.Context(), common.TXN_TOKEN_CONTEXT_KEY, txnToken)

	var orderRequest OrderRequest

	username := r.Header.Get("alpha-stock-user-name")
	if username == "" {
		h.Logger.Error("Unable to extract username from the header of the order request.")
		http.Error(w, "Unable to extract username from the header", http.StatusInternalServerError)

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		h.Logger.Error("Failed to decode order request body.", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	orderDetails, err := h.Service.Order(ctx, username, orderRequest.StockID, orderRequest.OrderType, orderRequest.Quantity)
	if err != nil {
		h.Logger.Error("Failed to process stock order.", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(orderDetails); err != nil {
		h.Logger.Error("Failed to encode response of a stock-order request.", zap.Error(err))

		return
	}

	h.Logger.Info("Order request processed successfully.")
}
