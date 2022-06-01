package main

import (
	"log"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/global/config"
	"github.com/nmluci/sumber-sari-garden/internal/global/middleware"
	"github.com/nmluci/sumber-sari-garden/internal/global/router"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/httputil"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("[main] starting server")
	// config.Init() //! Uncomment this to local dev

	conf := config.GetConfig()
	root := mux.NewRouter()
	db := database.Init()

	arrayWhitelistedUrls := strings.Split(conf.WhitelistUrl, ",")
	mapWhitelistedUrls := make(map[string]bool)

	for _, v := range arrayWhitelistedUrls {
		mapWhitelistedUrls[v] = true
	}

	root.Use(middleware.ErrorHandlingMiddleware)
	root.Use(middleware.CorsMiddleware(mapWhitelistedUrls))
	root.Use(middleware.ContentTypeMiddleware)

	router.Init(root, db)

	httputil.ProvideServer(conf.ServerAddress, root).ListerAndServe()
	log.Println("[main] stopping server gracefully")
}
