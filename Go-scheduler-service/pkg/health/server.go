package health

import (
	"Go-scheduler-service/internal/domain/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Srv         *http.Server
	healthPort  int
	pingFreq    int
	pingHost    string
	pingPort    int
	pingUrl     string
	stopChannel chan os.Signal
	logger      logger.Logger
}

func New(healthPort int, pingHost string, pingPort int, pingUrl string, freq int, logger logger.Logger, stopChannel chan os.Signal) (*Server, error) {
	srv := &http.Server{
		Addr: ":" + strconv.Itoa(healthPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(fmt.Sprintf("Failed to start Health Server [after calling ListenAndServe]: %s", err.Error()))
		}
	}()

	return &Server{
		Srv:         srv,
		healthPort:  healthPort,
		pingFreq:    freq,
		pingHost:    pingHost,
		pingPort:    pingPort,
		pingUrl:     pingUrl,
		stopChannel: stopChannel,
		logger:      logger,
	}, nil
}

func (hs *Server) HealthCheck() {
	pingTimer := time.NewTicker(time.Duration(hs.pingFreq) * time.Minute)
	defer pingTimer.Stop()

	for {
		select {
		case <-pingTimer.C:
			url := "http://" + hs.pingHost + ":" + strconv.Itoa(hs.pingPort) + hs.pingUrl
			hs.logger.Info(fmt.Sprintf("Ping %s ...", url))

			client := &http.Client{Timeout: 5 * time.Second}

			resp, err := client.Get(url)

			if err == nil && resp.StatusCode == http.StatusOK {
				res, err := io.ReadAll(resp.Body)
				if err != nil {
					hs.logger.Warning(fmt.Sprintf("Failed to decode response: %s", err.Error()))
				} else {
					hs.logger.Info(fmt.Sprintf("Received response: %v", string(res)))
				}
				hs.logger.Info("Ping successful.")
			} else {
				hs.logger.Warning("Unsuccessful ping.")
				hs.stopChannel <- os.Interrupt
			}
		}
	}
}
