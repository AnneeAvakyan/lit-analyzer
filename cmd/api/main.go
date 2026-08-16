package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AnneeAvakyan/litanalyzer/internal/config"
	"github.com/AnneeAvakyan/litanalyzer/internal/repository/postgres"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}
	log.Println("connected to database")

	bookRepo := postgres.NewPostgresBookRepository(pool)
	characterRepo := postgres.NewPostgresCharacterRepository(pool)
	chapterRepo := postgres.NewPostgresChapterRepository(pool)
	mentionRepo := postgres.NewPostgresMentionRepository(pool)
	aliasRepo := postgres.NewPostgresAliasRepository(pool)
	relationshipRepo := postgres.NewPostgresRelationshipRepository(pool)

	// пока просто чтобы использовать переменные и не ловить "unused variable"
	_ = bookRepo
	_ = characterRepo
	_ = chapterRepo
	_ = mentionRepo
	_ = aliasRepo
	_ = relationshipRepo

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Printf("starting server on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
