package main

import (
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/getenv"
	"github.com/nmluci/sumber-sari-garden/pkg/router/controller"
	globalRouter "github.com/nmluci/sumber-sari-garden/pkg/router/global"
	"github.com/nmluci/sumber-sari-garden/pkg/server"
)

func main() {
	envVariable := getenv.GetEnvironmentVariable()
	db := database.InitDatabaseFromEnvVariable(envVariable)
	r := globalRouter.InitGlobalRouter(envVariable["WHITELISTED_URLS"])
	controller.InitController(r, db)
	s := server.ProvideServer(envVariable["SERVER_ADDRESS"], r)
	s.ListerAndServe()
}