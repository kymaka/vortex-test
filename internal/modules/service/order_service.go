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

/*
GetOrderBook retrieves the order book for a specific exchange and trading pair.
Converts the order book model to a DTO before returning.
*/
func (osi *orderServiceImpl) GetOrderBook(exchangeName, pair string) (*models.OrderBookDTO, error) {
	order, err := osi.repo.FindOrder(exchangeName, pair)
	if err != nil {
		return nil, err
	}

	orderDTO := order.ToDTO()
	return &orderDTO, nil
}

/*
SaveOrderBook saves the order book details using the provided parameters.
Converts the DTO to a model before saving to the repository.
*/
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

// GetOrderHistory retrieves the order history for a given client.
func (osi *orderServiceImpl) GetOrderHistory(client *models.Client) ([]*models.HistoryOrder, error) {
	return osi.repo.FindOrderHistory(client)
}

/*
SaveOrder saves an order history record for a given client.
Adds client details to the order before saving to the repository.
*/
func (osi *orderServiceImpl) SaveOrder(client *models.Client, order *models.HistoryOrder) error {
	newOrder := *order
	newOrder.ClientName = client.ClientName
	newOrder.ExchangeName = client.ExchangeName
	newOrder.Label = client.Label
	newOrder.Pair = client.Pair

	return osi.repo.SaveOrderHistory(newOrder)
}
