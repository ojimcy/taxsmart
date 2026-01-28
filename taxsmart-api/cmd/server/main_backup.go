package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/taxsmart/taxsmart-api/internal/config"
	"github.com/taxsmart/taxsmart-api/internal/handler"
	"github.com/taxsmart/taxsmart-api/internal/middleware"
)

func main() {
	// Load .env file if present
	godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Create handlers
	h := handler.NewHandler(cfg.AIProvider, cfg.AIAPIKey)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/health", h.HealthCheck)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public endpoints (no auth required for parsing/classification)
		r.Post("/parse", h.ParseFile)
		r.Post("/classify", h.ClassifyTransactions)
		r.Post("/tax/quick-pit", h.QuickCalculatePIT)

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.SupabaseJWTSecret))
			r.Post("/tax/calculate", h.CalculateTax)
		})
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 TaxSmart API starting on port %s", port)
	log.Printf("📊 AI Provider: %s", cfg.AIProvider)
	if cfg.AIAPIKey != "" {
		log.Printf("🤖 AI Classification: Enabled")
	} else {
		log.Printf("📋 AI Classification: Disabled (rule-based only)")
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
