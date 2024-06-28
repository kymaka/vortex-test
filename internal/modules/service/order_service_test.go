package service

import (
	"errors"
	"testing"

	"github.com/kymaka/vortex-test/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of the OrderRepository interface.
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) FindOrder(exchangeName, pair string) ([]*models.OrderBook, error) {
	args := m.Called(exchangeName, pair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.OrderBook), args.Error(1)
}

func (m *MockOrderRepository) SaveOrder(order models.OrderBook) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindOrderHistory(client *models.Client) ([]*models.HistoryOrder, error) {
	args := m.Called(client)
	return args.Get(0).([]*models.HistoryOrder), args.Error(1)
}

func (m *MockOrderRepository) SaveOrderHistory(order models.HistoryOrder) error {
	args := m.Called(order)
	return args.Error(0)
}

func TestGetOrderBook(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	exchangeName := "test_exchange"
	pair := "BTC/USD"
	order := []*models.OrderBook{
		{
			Exchange: exchangeName,
			Pair:     pair,
		},
	}

	mockRepo.On("FindOrder", exchangeName, pair).Return(order, nil)

	result, err := service.GetOrderBook(exchangeName, pair)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exchangeName, result[0].Exchange)
	assert.Equal(t, pair, result[0].Pair)

	mockRepo.AssertExpectations(t)
}

func TestGetOrderBook_Error(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	exchangeName := "test_exchange"
	pair := "BTC/USD"

	mockRepo.On("FindOrder", exchangeName, pair).Return(nil, errors.New("order not found"))

	result, err := service.GetOrderBook(exchangeName, pair)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSaveOrderBook(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	id := int64(1)
	exchangeName := "test_exchange"
	pair := "BTC/USD"
	asks := []*models.DepthOrder{}
	bids := []*models.DepthOrder{}

	orderDTO := models.OrderBookDTO{
		ID:       id,
		Exchange: exchangeName,
		Pair:     pair,
		Asks:     asks,
		Bids:     bids,
	}

	order := orderDTO.ToOrderBook()

	mockRepo.On("SaveOrder", order).Return(nil)

	err := service.SaveOrderBook(id, exchangeName, pair, asks, bids)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetOrderHistory(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	client := &models.Client{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "BTC/USD",
	}

	orderHistory := []*models.HistoryOrder{
		{
			ClientName:   client.ClientName,
			ExchangeName: client.ExchangeName,
			Label:        client.Label,
			Pair:         client.Pair,
		},
	}

	mockRepo.On("FindOrderHistory", client).Return(orderHistory, nil)

	result, err := service.GetOrderHistory(client)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderHistory, result)

	mockRepo.AssertExpectations(t)
}

func TestSaveOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	client := &models.Client{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "BTC/USD",
	}

	order := &models.HistoryOrder{}

	expectedOrder := models.HistoryOrder{
		ClientName:   client.ClientName,
		ExchangeName: client.ExchangeName,
		Label:        client.Label,
		Pair:         client.Pair,
	}

	mockRepo.On("SaveOrderHistory", expectedOrder).Return(nil)

	err := service.SaveOrder(client, order)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
