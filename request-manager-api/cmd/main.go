package main

import (
	"context"
	"os"
	"os/signal"
	"request_manager_api"
	"request_manager_api/pkg/handlers"
	"request_manager_api/pkg/repository"
	"request_manager_api/pkg/services"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := loadConfig()

	db, err := repository.NewMysqlDb(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db, cfg)
	services := services.NewService(repos)
	handlers := handlers.NewHandler(services)

	srv := new(request_manager_api.Server)
	port := "8000"

	go func() {
		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error running server: %s", err.Error())
		}
	}()
	logrus.Printf("E-RequestControl started on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("E-RequestControl shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func loadConfig() repository.Config {
	return repository.Config{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		Username: getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Dbname:   getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		logrus.Fatalf("environment variable %s not set", key)
	}
	return val
}
