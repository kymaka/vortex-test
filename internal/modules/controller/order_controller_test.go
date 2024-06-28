package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kymaka/vortex-test/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockOrderService struct{}

func (m *MockOrderService) GetOrderBook(exchangeName, pair string) ([]*models.OrderBookDTO, error) {
	if exchangeName == "invalid" || pair == "invalid" {
		return nil, gorm.ErrRecordNotFound
	}

	if exchangeName == "error" {
		return nil, gorm.ErrInvalidValue
	}

	return []*models.OrderBookDTO{
		{
			Exchange: exchangeName,
			Pair:     pair,
			Asks:     []*models.DepthOrder{},
			Bids:     []*models.DepthOrder{},
		},
	}, nil
}

func (m *MockOrderService) SaveOrderBook(id int64, exchange, pair string, asks, bids []*models.DepthOrder) error {
	if exchange == "error" || pair == "error" {
		return errors.New("error saving order book")
	}
	return nil
}

func (m *MockOrderService) GetOrderHistory(client *models.Client) ([]*models.HistoryOrder, error) {
	if client.ClientName == "notfound" {
		return nil, gorm.ErrRecordNotFound
	}

	if client.ClientName == "error" {
		return nil, gorm.ErrInvalidValue
	}

	return []*models.HistoryOrder{}, nil
}

func (m *MockOrderService) SaveOrder(client *models.Client, history *models.HistoryOrder) error {
	if client.ClientName == "error" || history.Type == "error" {
		return errors.New("error saving order")
	}
	return nil
}

func TestGetOrderBookHandler(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest("GET", "/order/book?exchangeName=test&pair=ETH-BTC", nil)
	rr := httptest.NewRecorder()

	controller.GetOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var orders []*models.OrderBook
	err := json.NewDecoder(rr.Body).Decode(&orders)
	assert.NoError(t, err)
	assert.Equal(t, "test", orders[0].Exchange)
	assert.Equal(t, "ETH-BTC", orders[0].Pair)
}

func TestSaveOrderBookHandler(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	order := models.OrderBookDTO{
		Exchange: "test",
		Pair:     "ETH-BTC",
		Asks:     []*models.DepthOrder{},
		Bids:     []*models.DepthOrder{},
	}
	body, _ := json.Marshal(order)

	req := httptest.NewRequest("POST", "/order/book", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.SaveOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetOrderHistoryHandler(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	client := models.Client{
		ClientName: "testclient",
	}
	body, _ := json.Marshal(client)

	req := httptest.NewRequest("GET", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.GetOrderHistoryHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var orders []models.HistoryOrder
	err := json.NewDecoder(rr.Body).Decode(&orders)
	assert.NoError(t, err)
	assert.Empty(t, orders)
}

func TestGetOrderHistoryHandler_RecordNotFound(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	client := models.Client{
		ClientName: "notfound",
	}
	body, _ := json.Marshal(client)

	req := httptest.NewRequest("GET", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.GetOrderHistoryHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetOrderHistoryHandler_BadRequest(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest("GET", "/order/history", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	controller.GetOrderHistoryHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveOrderHandler(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	payload := models.HistoryOrderPayload{
		Client: models.Client{
			ClientName: "testclient",
		},
		History: models.HistoryOrder{
			Type: "buy",
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.SaveOrderHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSaveOrderBookHandler_MissingRequiredFields(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest(
		"POST",
		"/order/book",
		bytes.NewReader([]byte(`{"exchange": "test", "asks": [], "bids": []}`)))
	rr := httptest.NewRecorder()

	controller.SaveOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveOrderBookHandler_InternalServerError(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	order := models.OrderBookDTO{
		Exchange: "error",
		Pair:     "error",
		Asks:     []*models.DepthOrder{},
		Bids:     []*models.DepthOrder{},
	}
	body, _ := json.Marshal(order)

	req := httptest.NewRequest("POST", "/order/book", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.SaveOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSaveOrderHandler_MissingRequiredFields(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	payload := models.HistoryOrderPayload{
		Client: models.Client{},
		History: models.HistoryOrder{
			Type: "buy",
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.SaveOrderHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveOrderHandler_InternalServerError(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	payload := models.HistoryOrderPayload{
		Client: models.Client{
			ClientName: "error",
		},
		History: models.HistoryOrder{
			Type: "error",
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.SaveOrderHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOrderHistoryHandler_InternalServerError(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	client := models.Client{
		ClientName: "error",
	}
	body, _ := json.Marshal(client)

	req := httptest.NewRequest("GET", "/order/history", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.GetOrderHistoryHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOrderHandler_BadRequest(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest("GET", "/order/book?exchangeName=&pair=ETH-BTC", nil)
	rr := httptest.NewRecorder()

	controller.GetOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetOrderHandler_NotFound(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest("GET", "/order/book?exchangeName=invalid&pair=ETH-BTC", nil)
	rr := httptest.NewRecorder()

	controller.GetOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetOrderHandler_Error(t *testing.T) {
	controller := NewOrderController(&MockOrderService{})

	req := httptest.NewRequest("GET", "/order/book?exchangeName=error&pair=ETH-BTC", nil)
	rr := httptest.NewRecorder()

	controller.GetOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
