package main

import (
	"Go-scheduler-service/pkg/health"
	"Go-scheduler-service/pkg/logger"
	scheduler2 "Go-scheduler-service/pkg/scheduler"
	"context"
	"fmt"
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

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("[%s] pong", time.Now())))
		if err != nil {
			myLogger.Warning("Failed to write response")
		}
	})

	defer cancel()

	scheduler := scheduler2.New("localhost", 8080, "/user", 1, 5, myLogger)
	go func() {
		scheduler.ScheduleUsers()
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
