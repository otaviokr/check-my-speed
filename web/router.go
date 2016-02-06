package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kardianos/osext"
	log "gopkg.in/inconshreveable/log15.v2"
)

func NewRouter() *mux.Router {
	var handler http.Handler

	logger := log.New("module", "web.router")
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name, logger)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	currentFolder, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	handler = http.FileServer(http.Dir(currentFolder + "/html/lib/"))
	handler = http.StripPrefix("/resources/", handler)
	handler = Logger(handler, "Resources", logger)
	router.
		Methods("GET").
		PathPrefix("/resources/").
		Name("Resources").
		Handler(handler)

	return router
}
