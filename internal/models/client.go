package models

type Client struct {
	ClientName   string `json:"clientName"`
	ExchangeName string `json:"exchangeName"`
	Label        string `json:"label"`
	Pair         string `json:"pair"`
}
