package main

import (
	"Golang-practice-2023/internal/transport/rest/handler"
	"Golang-practice-2023/internal/user/repository"
	"Golang-practice-2023/internal/user/service"
	"Golang-practice-2023/pkg/health"
	"Golang-practice-2023/pkg/logger"
	"Golang-practice-2023/pkg/migration"
	"Golang-practice-2023/pkg/pgconnect"
	"Golang-practice-2023/pkg/pubsub/nats/pub"
	"Golang-practice-2023/pkg/pubsub/nats/sub"
	"context"
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
	fmt.Println("It's work")

	envPath, envErr := filepath.Abs("dev.env")
	if envErr != nil {
		log.Fatal("Can't get environment file")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	zeroLogLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	myLogger, err := logger.New(os.Getenv("LOG_LEVEL"), &zeroLogLogger)
	if err != nil {
		fmt.Println("* " + err.Error())
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		myLogger.Fatal("Failed to get Postgresql port")
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
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		myLogger.Fatal(err.Error())
	} else {
		myLogger.Warning(err.Error())
	}

	publisher, err := pub.New(fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")), myLogger)
	if err != nil {
		myLogger.Warning(err.Error())
	}

	subscriber, err := sub.New(fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")), myLogger)
	_, err = subscriber.Subscribe("NewUser", func(msg *nats.Msg) {
		fmt.Println("Received message: " + string(msg.Data))
	})

	userRepository := repository.New(db, myLogger)
	userService := service.New(userRepository, myLogger, publisher)
	userHandler := handler.New(userService, myLogger)

	if err != nil {
		myLogger.Warning("Failed to subscribe")
	}

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	userHandler.InitRoutes(router)
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("[%s] pong", time.Now())))
		if err != nil {
			// todo
		}
	})

	defer cancel()
	defer publisher.Conn.Close()

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

	fmt.Println("Hi! 1")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	fmt.Println("Hi! 2")

	healthSrv, err := health.New(os.Getenv("HEALTH_PORT"), 1, myLogger, c)
	if err != nil {
		//logger todo
		fmt.Println(err)
	}

	go func() {
		healthSrv.HealthCheck()
	}()

	fmt.Println("Hi! 3")

	sig := <-c
	myLogger.Info(fmt.Sprintf("shutting down the server, received signal : %s", sig.String()))
	ctx2, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx2)
	if err != nil {
		myLogger.Fatal("Could not shutdown the server (after getting signal): " + err.Error())
	}
}
