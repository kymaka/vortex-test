package db

import (
	"os"

	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

/*
Function for connecting to default db in ClickHouse
Uses .env variable for name, user, passm host and port
*/
func Connect(env string) (*gorm.DB, error) {
	godotenv.Load(env)

	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s", user, pass, host, port, name)

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/*
Migrations - creatin of tables if they not exists
Using raw sql queries because gorm currently doesnt't support PK for ClickHouse
*/
func Migrate(db *gorm.DB) error {
	if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS order_books (
				id Int64,
				exchange String,
				pair String,
				asks String,
				bids String
			) ENGINE = MergeTree()
			PRIMARY KEY (exchange, pair)
			ORDER BY (exchange, pair);`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS history_orders (
				client_name String,
				exchange_name String,
				label String,
				pair String,
				side String,
				type String,
				base_qty Float64,
				price Float64,
				algorithm_name_placed String,
				lowest_sell_prc Float64,
				highest_buy_prc Float64,
				commission_quote_qty Float64,
				time_placed DateTime
			) ENGINE = MergeTree()
			PRIMARY KEY (client_name, exchange_name, pair)
			ORDER BY (client_name, exchange_name, pair);`).Error; err != nil {
		return err
	}

	return nil
}
