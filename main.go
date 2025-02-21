package main

import (
	"log"
	"net/http"
	"time"

	"github.com/kymaka/vortex-test/internal/infrastructure/db"
	"github.com/kymaka/vortex-test/internal/modules/controller"
	"github.com/kymaka/vortex-test/internal/modules/repository"
	"github.com/kymaka/vortex-test/internal/modules/service"

	_ "github.com/kymaka/vortex-test/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
)

func main() {
	gormDB, err := db.Connect(".env")
	if err != nil {
		log.Fatalf("failed to connect to clickhouse: %v", err)
	}
	log.Println("Connected to ClickHouse successfully!")

	err = db.Migrate(gormDB)
	if err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}

	repository := repository.NewOrderRepository(gormDB)
	service := service.NewOrderService(repository)
	controller := controller.NewOrderController(service)

	r := chi.NewMux()

	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Second))

		r.Get("/order/book", controller.GetOrderBookHandler)
		r.Get("/order/history", controller.GetOrderHistoryHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(200, 1*time.Second))

		r.Post("/order/book", controller.SaveOrderBookHandler)
		r.Post("/order/history", controller.SaveOrderHandler)
	})

	http.ListenAndServe(":8080", r)
}
