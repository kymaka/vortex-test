package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"vortex_test/internal/models"
	"vortex_test/internal/modules/service"

	"gorm.io/gorm"
)

type OrderController interface {
	GetOrderBookHandler(w http.ResponseWriter, r *http.Request)
	SaveOrderBookHandler(w http.ResponseWriter, r *http.Request)
	GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request)
	SaveOrderHandler(w http.ResponseWriter, r *http.Request)
}

type orderControllerImpl struct {
	service service.OrderService
}

func NewOrderController(s service.OrderService) OrderController {
	return &orderControllerImpl{service: s}
}

func (oci *orderControllerImpl) GetOrderBookHandler(w http.ResponseWriter, r *http.Request) {
	exchangeName := r.URL.Query().Get("exchangeName")
	pair := r.URL.Query().Get("pair")

	if exchangeName == "" || pair == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := oci.service.GetOrderBook(exchangeName, pair)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(order)
	w.Write(bytes)

}

func (oci *orderControllerImpl) SaveOrderBookHandler(w http.ResponseWriter, r *http.Request) {
	var order models.OrderBookDTO
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil || order.Pair == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = oci.service.SaveOrderBook(order.ID, order.Exchange, order.Pair, order.Asks, order.Bids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (oci *orderControllerImpl) GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := oci.service.GetOrderHistory(&client)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(orders)
	w.Write(bytes)
}

func (oci *orderControllerImpl) SaveOrderHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.HistoryOrderPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.Client.ClientName == "" || payload.History.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = oci.service.SaveOrder(&payload.Client, &payload.History)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
