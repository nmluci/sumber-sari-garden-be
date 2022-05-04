package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/ping"
	"github.com/nmluci/sumber-sari-garden/pkg/middleware"
)

func InitController(globalRouter *mux.Router, db *sql.DB) {
	globalRouter.Use(middleware.ErrorHandlingMiddleware)

	router := globalRouter.NewRoute().Subrouter()

	pingService := ping.NewPingService()
	pingController := ping.NewPingHandler(router, pingService)
	pingController.InitHandler()
}
