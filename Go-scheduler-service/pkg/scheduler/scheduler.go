package scheduler

import (
	"Go-scheduler-service/internal/domain/logger"
	"Go-scheduler-service/internal/domain/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Scheduler struct {
	host   string
	port   int
	url    string
	freq   int
	limit  int
	logger logger.Logger
}

func New(host string, port int, url string, freq int, limit int, logger logger.Logger) *Scheduler {
	return &Scheduler{
		host:   host,
		port:   port,
		url:    url,
		freq:   freq,
		limit:  limit,
		logger: logger,
	}
}

func (s *Scheduler) ScheduleUsers() {
	pingTimer := time.NewTicker(time.Duration(s.freq) * time.Minute)
	defer pingTimer.Stop()

	offset := 0
	for {
		select {
		case <-pingTimer.C:
			users, err := s.getUsers(offset)
			if err != nil {
				s.logger.Warning(fmt.Sprintf("Error getting users in scheduler: %s", err.Error()))
				continue
			} else {
				offset += s.limit
			}

			if users != nil {
				s.logger.Info(fmt.Sprintf("%v", users))
			} else {
				s.logger.Info("No more users")
			}
		}
	}
}

func (s *Scheduler) getUsers(offset int) ([]user.User, error) {
	url := "http://" + s.host + ":" + strconv.Itoa(s.port) + s.url + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(s.limit)
	s.logger.Info(fmt.Sprintf("Send query to get users %s ...", url))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []user.User

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		s.logger.Warning(fmt.Sprintf("Error decoding users in scheduler: %s", err.Error()))
		return nil, err
	}

	return users, nil
}
