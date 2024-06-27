package service

import (
	"vortex_test/internal/models"
	"vortex_test/internal/modules/repository"
)

type OrderService interface {
	GetOrderBook(exchangeName, pair string) (*models.OrderBookDTO, error)
	SaveOrderBook(id int64, exchangeName, pair string, asks, bids []*models.DepthOrder) error
	GetOrderHistory(client *models.Client) ([]*models.HistoryOrder, error)
	SaveOrder(client *models.Client, order *models.HistoryOrder) error
}

type orderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderServiceImpl{repo: r}
}

func (osi *orderServiceImpl) GetOrderBook(exchangeName, pair string) (*models.OrderBookDTO, error) {
	order, err := osi.repo.FindOrder(exchangeName, pair)
	if err != nil {
		return nil, err
	}

	orderDTO := order.ToDTO()
	return &orderDTO, nil
}

func (osi *orderServiceImpl) SaveOrderBook(id int64, exchangeName, pair string, asks, bids []*models.DepthOrder) error {
	orderDTO := models.OrderBookDTO{
		ID:       id,
		Exchange: exchangeName,
		Pair:     pair,
		Asks:     asks,
		Bids:     bids,
	}

	order := orderDTO.ToOrderBook()
	return osi.repo.SaveOrder(order)
}

func (osi *orderServiceImpl) GetOrderHistory(client *models.Client) ([]*models.HistoryOrder, error) {
	return osi.repo.FindOrderHistory(client)
}

func (osi *orderServiceImpl) SaveOrder(client *models.Client, order *models.HistoryOrder) error {
	newOrder := *order
	newOrder.ClientName = client.ClientName
	newOrder.ExchangeName = client.ExchangeName
	newOrder.Label = client.Label
	newOrder.Pair = client.Pair

	return osi.repo.SaveOrderHistory(newOrder)
}
