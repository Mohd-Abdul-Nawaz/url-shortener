package main

import (
	"log"
	"net/http"
	"os"

	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
	"url-shortener/internal/store"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisStore, err := store.NewRedisStore(redisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisStore.Close()

	shortenerHandler := handler.NewShortenerHandler(redisStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenerHandler.ShortenURL)
	mux.HandleFunc("/", shortenerHandler.Redirect)

	// Apply middleware: logging, CORS, rate limiting
	loggedMux := middleware.Logging(mux)
	corsMux := middleware.CORS(loggedMux)
	rateLimited := middleware.RateLimit(corsMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting URL shortener on port %s", port)
	if err := http.ListenAndServe(":"+port, rateLimited); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}