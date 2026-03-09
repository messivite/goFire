package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/messivite/goFire/config"
	"github.com/messivite/goFire/server"
)

var h http.Handler

func init() {
	cfg := config.LoadFromEnv()

	var err error
	h, err = server.NewHandler(cfg)
	if err != nil {
		log.Printf("GoFire Vercel: handler init error: %v", err)
		return
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "server not configured: check environment variables",
		})
		return
	}

	if p := r.URL.Query().Get("__path"); p != "" {
		r.URL.Path = "/" + p
		q := r.URL.Query()
		q.Del("__path")
		r.URL.RawQuery = q.Encode()
	}
	h.ServeHTTP(w, r)
}
