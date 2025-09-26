package app

import (
	"context"
	"fmt"
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
	logger "docs_storage/pkg/logger"
	repository "docs_storage/internal/repository"
	storage "docs_storage/internal/storage"
)

type App struct {
	server *http.Server	
	config *Config
	logger *logger.Logger
}

func NewApp(config *Config) *App {
	l := logger.New(os.Stdout, os.Stderr)
	app := &App{
		config: config,
		logger: l,
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
		a.logger.Error.Println("Failed to initialize postgres:", err)
		return err
	}
	defer postgres.Close()

	docsRepo := repository.NewDocsRepo(postgres.Pool)
	userRepo := repository.NewUserRepo(postgres.Pool)
	sessionRepo := repository.NewSessionRepo(postgres.Pool)

	fileStorage := storage.NewLocalFileStorage("/app/files")

	docsSvc := service.NewDocsService(docsRepo, fileStorage, sessionRepo)
	authSvc := service.NewAuthService(userRepo, sessionRepo, a.config.Admin.token)

	docsHandler := handlers.NewDocsHandler(docsSvc, a.logger)
	authHandler := handlers.NewAuthHandler(authSvc, a.logger)
	
	router := mux.NewRouter()

	routes.SetupDocsRoutes(router, docsHandler)
	routes.SetupAuthRoutes(router, authHandler)
	
	serverAddr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)

	a.server = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		a.logger.Info.Printf("Starting server on %s", serverAddr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error.Fatalf("Error starting server: %v", err)
		}
	}()

	<-quit
	a.logger.Info.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(ctx, time.Duration(a.config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error.Println("Server forced to shutdown:", err)
		return err
	}

	a.logger.Info.Println("Server exited properly")
	return nil
}
