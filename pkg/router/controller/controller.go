package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	pingController "github.com/nmluci/sumber-sari-garden/internal/ping/controller"
	pingServicePkg "github.com/nmluci/sumber-sari-garden/internal/ping/service/impl"
	"github.com/nmluci/sumber-sari-garden/pkg/middleware"
)

func InitController(globalRouter *mux.Router, db *sql.DB) {
	globalRouter.Use(middleware.ErrorHandlingMiddleware)

	router := globalRouter.NewRoute().Subrouter()

	pingService := pingServicePkg.ProvidePingService()
	pingController := pingController.ProvideMsibController(router, pingService)
	pingController.InitController()
}