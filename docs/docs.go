// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/order/book": {
            "get": {
                "description": "Returns the order books for a given exchange and pair.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Get order books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Exchange Name",
                        "name": "exchangeName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Trading Pair",
                        "name": "pair",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.OrderBook"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Saves the order book details for a given exchange and pair.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Save order book",
                "parameters": [
                    {
                        "description": "Order Book DTO",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.OrderBookDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/order/history": {
            "get": {
                "description": "Returns the order history for a given client.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Get order history",
                "parameters": [
                    {
                        "description": "Client",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Client"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.HistoryOrder"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Saves an order for a given client.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Save order",
                "parameters": [
                    {
                        "description": "History Order Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.HistoryOrderPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Client": {
            "type": "object",
            "properties": {
                "clientName": {
                    "type": "string"
                },
                "exchangeName": {
                    "type": "string"
                },
                "label": {
                    "type": "string"
                },
                "pair": {
                    "type": "string"
                }
            }
        },
        "models.DepthOrder": {
            "type": "object",
            "properties": {
                "baseQty": {
                    "type": "number"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "models.HistoryOrder": {
            "type": "object",
            "properties": {
                "algorithmNamePlaced": {
                    "type": "string"
                },
                "baseQty": {
                    "type": "number"
                },
                "clientName": {
                    "type": "string"
                },
                "commissionQuoteQty": {
                    "type": "number"
                },
                "exchangeName": {
                    "type": "string"
                },
                "highestBuyPrc": {
                    "type": "number"
                },
                "label": {
                    "type": "string"
                },
                "lowestSellPrc": {
                    "type": "number"
                },
                "pair": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "side": {
                    "type": "string"
                },
                "timePlaced": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models.HistoryOrderPayload": {
            "type": "object",
            "properties": {
                "client": {
                    "$ref": "#/definitions/models.Client"
                },
                "history": {
                    "$ref": "#/definitions/models.HistoryOrder"
                }
            }
        },
        "models.OrderBook": {
            "type": "object",
            "properties": {
                "asks": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "bids": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "exchange": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "pair": {
                    "type": "string"
                }
            }
        },
        "models.OrderBookDTO": {
            "type": "object",
            "properties": {
                "asks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.DepthOrder"
                    }
                },
                "bids": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.DepthOrder"
                    }
                },
                "exchange": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "pair": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
