package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type OrderBookDTO struct {
	ID       int64
	Exchange string
	Pair     string
	Asks     []*DepthOrder
	Bids     []*DepthOrder
}

type Tuple [2]float64

type Tuples []Tuple

func (t Tuples) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Tuples) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return errors.New("unsupported type for Tuples")
	}
}

type OrderBook struct {
	ID       int64
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
	Asks     Tuples `json:"asks"`
	Bids     Tuples `json:"bids"`
}

func (dto *OrderBookDTO) ToOrderBook() OrderBook {
	return OrderBook{
		ID:       dto.ID,
		Exchange: dto.Exchange,
		Pair:     dto.Pair,
		Asks:     depthOrdersToTuples(dto.Asks),
		Bids:     depthOrdersToTuples(dto.Bids),
	}
}

func (o *OrderBook) ToDTO() OrderBookDTO {
	return OrderBookDTO{
		ID:       o.ID,
		Exchange: o.Exchange,
		Pair:     o.Pair,
		Asks:     tuplesToDepthOrders(o.Asks),
		Bids:     tuplesToDepthOrders(o.Bids),
	}
}

func depthOrdersToTuples(orders []*DepthOrder) Tuples {
	tuples := make(Tuples, len(orders))
	for i, order := range orders {
		tuples[i] = Tuple{order.Price, order.BaseQty}
	}
	return tuples
}

func tuplesToDepthOrders(tuples Tuples) []*DepthOrder {
	orders := make([]*DepthOrder, len(tuples))
	for i, tuple := range tuples {
		orders[i] = &DepthOrder{Price: tuple[0], BaseQty: tuple[1]}
	}
	return orders
}
