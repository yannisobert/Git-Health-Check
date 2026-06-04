package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	port string
}

func New() *Server {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{port: port}
}

func (s *Server) Run() error {
	h := newHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/analyze", h.analyze)
	mux.HandleFunc("/api/history", h.history)
	mux.HandleFunc("/api/compare", h.compare)
	mux.HandleFunc("/api/rivals", h.rivals)
	mux.Handle("/", spaHandler())

	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("ghhealth listening on http://localhost%s", addr)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
