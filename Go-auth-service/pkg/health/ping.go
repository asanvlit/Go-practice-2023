package health

import (
	"Golang-practice-2023/internal/domain/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Srv    *http.Server
	freq   int32
	logger logger.Logger
}

func New(port string, freq int32, logger logger.Logger) (*Server, error) {
	srv := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	return &Server{
		Srv:    srv,
		freq:   freq,
		logger: logger,
	}, nil
}

func (hs *Server) HealthCheck() {
	go func() {
		if err := hs.Srv.ListenAndServe(); err != nil {
			hs.logger.Fatal(fmt.Sprintf("Health check server error [after calling ListenAndServe]: " + err.Error()))
		}
	}()

	for {
		time.Sleep(time.Duration(hs.freq) * time.Second)

		resp, err := http.Get("http://" + os.Getenv("HOST") + ":" + os.Getenv("PORT"))
		if err != nil || resp.StatusCode != http.StatusOK {
			("Server is down, shutting down gracefully...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				fmt.Println("Could not shutdown the server:", err)
				os.Exit(1)
			}

			break
		}
	}
}
