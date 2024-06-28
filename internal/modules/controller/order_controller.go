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

// GetOrderBookHandler retrieves the order book for a specific exchange and pair.
//
//	@Summary		Get order book
//	@Description	Returns the order book for a given exchange and pair.
//	@Tags			orders
//	@Produce		json
//	@Param			exchangeName	query		string	true	"Exchange Name"
//	@Param			pair			query		string	true	"Trading Pair"
//	@Success		200				{object}	models.OrderBook
//	@Failure		400				{string}	string	"Bad Request"
//	@Failure		404				{string}	string	"Not Found"
//	@Failure		500				{string}	string	"Internal Server Error"
//	@Router			/order/book [get]
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

// SaveOrderBookHandler saves the order book details.
//
//	@Summary		Save order book
//	@Description	Saves the order book details for a given exchange and pair.
//	@Tags			orders
//	@Accept			json
//	@Param			order	body		models.OrderBookDTO	true	"Order Book DTO"
//	@Success		200		{string}	string				"OK"
//	@Failure		400		{string}	string				"Bad Request"
//	@Failure		500		{string}	string				"Internal Server Error"
//	@Router			/order/book [post]
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

// GetOrderHistoryHandler retrieves the order history for a client.
//
//	@Summary		Get order history
//	@Description	Returns the order history for a given client.
//	@Tags			orders
//	@Produce		json
//	@Param			client	body		models.Client	true	"Client"
//	@Success		200		{array}		models.HistoryOrder
//	@Failure		400		{string}	string	"Bad Request"
//	@Failure		404		{string}	string	"Not Found"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/order/history [post]
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

// SaveOrderHandler saves an order for a client.
//
//	@Summary		Save order
//	@Description	Saves an order for a given client.
//	@Tags			orders
//	@Accept			json
//	@Param			payload	body		models.HistoryOrderPayload	true	"History Order Payload"
//	@Success		200		{string}	string						"OK"
//	@Failure		400		{string}	string						"Bad Request"
//	@Failure		500		{string}	string						"Internal Server Error"
//	@Router			/order/history [post]
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
