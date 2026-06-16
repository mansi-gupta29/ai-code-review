package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mansi-gupta29/ai-code-review/internal/ai"
	"github.com/mansi-gupta29/ai-code-review/internal/handlers"
	"github.com/mansi-gupta29/ai-code-review/internal/store"
)

func main() {
	// 1. Connect to Postgres
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/codereview"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("db ping failed:", err)
	}
	log.Println("connected to postgres")

	// 2. Create the store (we'll use it next step)
	s := store.New(pool)
	aiClient := ai.New()

	h := handlers.New(s, aiClient)

	// 3. Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/reviews/{id}", h.GetReviewByID)
	r.Post("/review", h.CreateReview)

	// 4. Start server
	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
