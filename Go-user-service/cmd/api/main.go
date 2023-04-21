package main

import (
	"Golang-practice-2023/internal/domain/user"
	"Golang-practice-2023/internal/user/repository"
	"Golang-practice-2023/internal/user/service"
	"Golang-practice-2023/pkg/health"
	"Golang-practice-2023/pkg/logger"
	"Golang-practice-2023/pkg/migration"
	"Golang-practice-2023/pkg/pgconnect"
	"Golang-practice-2023/pkg/pubsub/nats/sub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
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

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

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

	ctx := context.Background()

	subscriber, err := sub.New(fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")), myLogger)
	_, err = subscriber.Subscribe("NewUser", func(msg *nats.Msg) {
		fmt.Println("Received message: " + string(msg.Data))

		var newUser *user.User
		err := json.Unmarshal(msg.Data, &newUser)
		if err != nil {
			myLogger.Warning("Failed to unmarshal user")
		} else {
			err := userService.Save(ctx, newUser)
			if err != nil {
				myLogger.Warning("Failed to create user " + err.Error())
			}
		}
	})

	if err != nil {
		myLogger.Warning("Failed to subscribe")
	}

	defer cancel()
	defer subscriber.Conn.Close()

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("[%s] pong", time.Now())))
		if err != nil {
			myLogger.Warning("Failed to write response")
		}
	})

	port := os.Getenv("PORT")
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
	healthSrv, err := health.New(healthPort, os.Getenv("HOST"), "/ping", 5, myLogger, c)
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
