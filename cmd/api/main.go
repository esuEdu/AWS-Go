package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/esuEdu/AWS-Go/internal/config"
	"github.com/esuEdu/AWS-Go/internal/handlers"
	"github.com/esuEdu/AWS-Go/internal/repository/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	if err := postgres.Migrate(ctx, pool); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.New(pool)
	r.Route("/items", func(r chi.Router) {
		r.Get("/", h.ListItems)
		r.Post("/", h.CreateItem)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetItem)
			r.Put("/", h.UpdateItem)
			r.Delete("/", h.DeleteItem)
		})
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Printf("server listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
