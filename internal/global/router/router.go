package router

import (
	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/auth"
	"github.com/nmluci/sumber-sari-garden/internal/ping"
	"github.com/nmluci/sumber-sari-garden/internal/product"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

func Init(globalRouter *mux.Router, db *database.DatabaseClient) {
	// globalRouter.Use(middleware.ErrorHandlingMiddleware)
	router := globalRouter.NewRoute().PathPrefix("/v1").Subrouter()

	pingService := ping.NewPingService()
	pingController := ping.NewPingHandler(router, pingService)
	pingController.InitHandler()

	authService := auth.NewAuthService(db)
	authController := auth.NewAuthHandler(router, authService)
	authController.InitHandler()

	productService := product.NewProductService(db)
	productController := product.NewProductHandler(router, productService)
	productController.InitHandler()
}
