package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"urlshort/internal/handlers"
	"urlshort/internal/logger"
	"urlshort/internal/storage"
	"urlshort/internal/storage/memory"
	"urlshort/internal/storage/postgres"
	"urlshort/internal/config"
)

func main() {
	cfg := config.New()

	log, err := logger.New()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer log.Sync()

	var store storage.Storage
	switch cfg.StorageType {
	case "memory":
		store = memory.NewMemoryStorage()
		log.Info("using in-memory storage")
	case "postgres":
		store, err = postgres.NewPostgresStorage(cfg.PostgresConn)
		if err != nil {
			log.Fatalw("failed to connect to postgres", "error", err)
		}
		log.Info("using postgres storage")
	default:
		log.Fatalw("unknown storage type", "type", cfg.StorageType)
	}
	defer store.Close()

	handler := handlers.NewHandler(store, log, cfg.BaseURL)

	r := mux.NewRouter()
	
	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/shorten", handler.Shorten).Methods("POST")
	
	r.HandleFunc("/api/save", handler.Save).Methods("POST")
	r.HandleFunc("/{short}", handler.Get).Methods("GET")

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	go func() {
		log.Infow("starting server", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")
}