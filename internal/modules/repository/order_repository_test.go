package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"vortex_test/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stretchr/testify/assert"
	c "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	// Connect to the ClickHouse server
	conn, err := sql.Open("clickhouse", "clickhouse://default:@localhost:9000")
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	// Ping the server to check the connection
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Printf("Ping failed: %v\n", err)
		}
		return nil, err
	}

	// Create a new database
	query := "CREATE DATABASE IF NOT EXISTS orders_test"
	_, err = conn.ExecContext(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	db, err := gorm.Open(c.Open("clickhouse://default:@localhost:9000/orders_test"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&models.HistoryOrder{}, &models.OrderBook{})
	return db, nil
}

func teardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw DB connection: %v", err)
	}
	sqlDB.Exec("DROP DATABASE orders_test")
	sqlDB.Close()
}

func TestFindOrder(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer teardownTestDB(db)

	repo := NewOrderRepository(db)
	exchangeName := "test_exchange"
	pair := "BTC/USD"

	order := models.OrderBook{Exchange: exchangeName, Pair: pair}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}

	foundOrder, err := repo.FindOrder(exchangeName, pair)
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	assert.Equal(t, exchangeName, foundOrder[0].Exchange)
	assert.Equal(t, pair, foundOrder[0].Pair)
}

func TestSaveOrder(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer teardownTestDB(db)

	repo := NewOrderRepository(db)
	order := models.OrderBook{Exchange: "test_exchange", Pair: "BTC/USD"}

	err = repo.SaveOrder(order)
	assert.NoError(t, err)

	var savedOrder models.OrderBook
	if err := db.First(&savedOrder, "exchange = ? AND pair = ?", order.Exchange, order.Pair).Error; err != nil {
		t.Fatalf("failed to find saved order: %v", err)
	}

	assert.Equal(t, order.Exchange, savedOrder.Exchange)
	assert.Equal(t, order.Pair, savedOrder.Pair)
}

func TestFindOrderHistory(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer teardownTestDB(db)

	repo := NewOrderRepository(db)
	client := &models.Client{ClientName: "test_client", ExchangeName: "test_exchange", Label: "test_label", Pair: "BTC/USD"}

	orderHistory := models.HistoryOrder{
		ClientName:   client.ClientName,
		ExchangeName: client.ExchangeName,
		Label:        client.Label,
		Pair:         client.Pair,
	}
	if err := db.Create(&orderHistory).Error; err != nil {
		t.Fatalf("failed to create test order history: %v", err)
	}

	foundHistory, err := repo.FindOrderHistory(client)
	assert.NoError(t, err)
	assert.NotNil(t, foundHistory)
	assert.Equal(t, client.ClientName, foundHistory[0].ClientName)
	assert.Equal(t, client.ExchangeName, foundHistory[0].ExchangeName)
	assert.Equal(t, client.Label, foundHistory[0].Label)
	assert.Equal(t, client.Pair, foundHistory[0].Pair)
}

func TestSaveOrderHistory(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer teardownTestDB(db)

	repo := NewOrderRepository(db)
	order := models.HistoryOrder{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "BTC/USD",
	}

	err = repo.SaveOrderHistory(order)
	assert.NoError(t, err)

	var savedOrderHistory models.HistoryOrder
	if err := db.First(
		&savedOrderHistory,
		"client_name = ? AND exchange_name = ? AND label = ? AND pair = ?",
		order.ClientName, order.ExchangeName, order.Label, order.Pair).Error; err != nil {
		t.Fatalf("failed to find saved order history: %v", err)
	}

	assert.Equal(t, order.ClientName, savedOrderHistory.ClientName)
	assert.Equal(t, order.ExchangeName, savedOrderHistory.ExchangeName)
	assert.Equal(t, order.Label, savedOrderHistory.Label)
	assert.Equal(t, order.Pair, savedOrderHistory.Pair)
}
