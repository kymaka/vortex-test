package repository

import (
	"vortex_test/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindOrder(exchangeName, pair string) (*models.OrderBook, error)
	SaveOrder(order models.OrderBook) error
	FindOrderHistory(client *models.Client) ([]*models.HistoryOrder, error)
	SaveOrderHistory(order models.HistoryOrder) error
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(d *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: d}
}

/*
FindOrder retrieves an order book from the database based on the exchange name and trading pair.
Returns the order book if found, or an error if not found or any other issue occurs.
*/
func (ori *orderRepositoryImpl) FindOrder(exchangeName, pair string) (*models.OrderBook, error) {
	var order models.OrderBook
	tx := ori.db.Where("exchange = ?", exchangeName).
		Where("pair = ?", pair).
		Find(&order)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &order, nil
}

/*
SaveOrder saves a new order book to the database.
Returns an error if the operation fails.
*/
func (ori *orderRepositoryImpl) SaveOrder(order models.OrderBook) error {
	tx := ori.db.Create(&order)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

/*
FindOrderHistory retrieves the order history for a given client from the database.
Returns the order history if found, or an error if not found or any other issue occurs.
*/
func (ori *orderRepositoryImpl) FindOrderHistory(client *models.Client) ([]*models.HistoryOrder, error) {
	var orderHistory []*models.HistoryOrder
	tx := ori.db.
		Where("client_name = ?", client.ClientName).
		Where("exchange_name = ?", client.ExchangeName).
		Where("label = ?", client.Label).
		Where("pair = ?", client.Pair).
		Find(&orderHistory)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return orderHistory, nil
}

/*
SaveOrderHistory saves a new order history record to the database.
Returns an error if the operation fails.
*/
func (ori *orderRepositoryImpl) SaveOrderHistory(order models.HistoryOrder) error {
	tx := ori.db.Create(&order)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
