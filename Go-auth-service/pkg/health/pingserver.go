package health

import (
	"Golang-practice-2023/internal/domain/logger"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Srv         *http.Server
	pingFreq    int32
	stopChannel chan os.Signal
	logger      logger.Logger // todo another params
}

func New(port string, freq int32, logger logger.Logger, stopChannel chan os.Signal) (*Server, error) {
	srv := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(fmt.Sprintf("Health check server error [after calling ListenAndServe]: " + err.Error()))
		}
	}()

	return &Server{
		Srv:         srv,
		pingFreq:    freq,
		stopChannel: stopChannel,
		logger:      logger,
	}, nil
}

func (hs *Server) HealthCheck() {
	hs.logger.Warning("Health server: started working")

	pingTimer := time.NewTicker(time.Duration(hs.pingFreq) * time.Minute)
	defer pingTimer.Stop()

	for {
		select {
		case <-pingTimer.C:
			hs.logger.Warning("Try ping...")

			client := &http.Client{Timeout: time.Second}

			resp, err := client.Get("http://" + os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/ping")

			if err != nil || resp.StatusCode != http.StatusOK {
				hs.logger.Warning("Error calling ping")
				hs.stopChannel <- os.Interrupt
			}

			hs.logger.Warning("Ping successful")
		}
	}
}
