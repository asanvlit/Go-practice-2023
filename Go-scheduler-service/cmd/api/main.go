package main

import (
	"Go-scheduler-service/internal/user/repository"
	"Go-scheduler-service/internal/user/service"
	"Go-scheduler-service/pkg/health"
	"Go-scheduler-service/pkg/logger"
	"Go-scheduler-service/pkg/migration"
	"Go-scheduler-service/pkg/pgconnect"
	scheduler2 "Go-scheduler-service/pkg/scheduler"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	envPath, envErr := filepath.Abs("dev.env")
	if envErr != nil {
		log.Fatal(fmt.Sprintf("Can't get environment file: %s", envErr))
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %s", err.Error()))
	}

	zeroLogLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	myLogger, err := logger.New(os.Getenv("LOG_LEVEL"), &zeroLogLogger)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating logger: %s", err))
	}

	ctx := context.Background()

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("[%s] Pong!", time.Now())))
		if err != nil {
			myLogger.Warning("Failed to write response")
		}
	})

	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		myLogger.Fatal(fmt.Sprintf("Failed to get Postgresql port: %s", err))
	}
	db, err := pgconnect.ConnectDatabase(pgconnect.ConnectionConfigData{
		Username:     os.Getenv("POSTGRES_USERNAME"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("POSTGRES_DATABASE"),
		Port:         pgPort,
		Host:         os.Getenv("POSTGRES_HOST"),
	})
	_ = db

	migrationData, err := migration.New(db.DB, "file://schemas")
	err = migrationData.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			myLogger.Warning(fmt.Sprintf("Did not migrate: %s", err.Error()))
		} else {
			myLogger.Fatal(fmt.Sprintf("Failed to migrate: %s", err.Error()))
		}
	}

	userRepository := repository.New(db, myLogger)
	userService := service.New(userRepository, myLogger)

	scheduler := scheduler2.New(userService, "go-auth", 8080, "/user", 1, 5, myLogger)
	go func() {
		scheduler.ScheduleUsers(ctx)
	}()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			return
		}
		if err != nil {
			myLogger.Fatal("Could not start the server (after http.ListenAndServe): " + err.Error())
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	healthPort, err := strconv.Atoi(os.Getenv("HEALTH_PORT"))
	if err != nil {
		myLogger.Fatal(fmt.Sprintf("Failed to get health port: %s", err.Error()))
	}

	pingPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		myLogger.Error("Failed to get PORT")
	}
	healthSrv, err := health.New(healthPort, os.Getenv("HOST"), pingPort, "/ping", 1, myLogger, c)
	if err != nil {
		myLogger.Fatal(fmt.Sprintf("Failed to start Health server: %s", err))
	}

	go func() {
		healthSrv.HealthCheck()
	}()

	sig := <-c
	myLogger.Info(fmt.Sprintf("shutting down the server, received signal : %s", sig.String()))
	ctx2, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx2)
	if err != nil {
		myLogger.Fatal("Could not shutdown the server (after getting signal): " + err.Error())
	}
}
