package web

import (
	"net/http"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

func Logger(inner http.Handler, name string, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		//fmt.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
		logger.Info("Request", "method", r.Method, "uri", r.RequestURI, "name", name, "duration", time.Since(start))
	})
}
