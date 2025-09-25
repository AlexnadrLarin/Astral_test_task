package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"github.com/gorilla/mux"

	handlers "docs_storage/internal/delivery/http/handlers"
	routes "docs_storage/internal/delivery/http/routes"
	service "docs_storage/internal/service"
	db "docs_storage/pkg/db"
	repository "docs_storage/internal/repository"
	storage "docs_storage/internal/storage"
)

type App struct {
	server *http.Server	
	config *Config
}

func NewApp(config *Config) *App {
	app := &App{
		config: config,
	}

	return app
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgres, err := db.NewPostgresWithConfig(db.PostgresConfig{
		Host:     a.config.Postgres.Host,
		Port:     a.config.Postgres.Port,
		Username: a.config.Postgres.Username,
		Password: a.config.Postgres.Password,
		DBName:   a.config.Postgres.DBName,
		SSLMode:  a.config.Postgres.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize postgres: %w", err)
	}
	defer postgres.Close()

	docsRepo := repository.NewDocsRepo(postgres.Pool)
	fileStorage := storage.NewLocalFileStorage("/app/files")

	docsSvc := service.NewDocsService(docsRepo, fileStorage)

	docsHandler := handlers.NewDocsHandler(docsSvc)
	
	router := mux.NewRouter()

	routes.SetupRoutes(router, docsHandler)
	
	serverAddr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)

	a.server = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(ctx, time.Duration(a.config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited properly")
	return nil
}
