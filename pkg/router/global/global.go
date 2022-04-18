package global

import (
	"strings"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/pkg/middleware"
)

func InitGlobalRouter(whitelistedUrls string) *mux.Router {
	r := mux.NewRouter()

	arrayWhitelistedUrls := strings.Split(whitelistedUrls, ",")

	mapWhitelistedUrls := make(map[string]bool)

	for _, v := range arrayWhitelistedUrls {
		mapWhitelistedUrls[v] = true
	}

	r.Use(middleware.ContentTypeMiddleware)
	r.Use(middleware.CorsMiddleware(mapWhitelistedUrls))
	return r
}