package scheduler

import (
	"Go-scheduler-service/internal/domain/apperrors"
	"Go-scheduler-service/internal/domain/logger"
	"Go-scheduler-service/internal/domain/user"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Scheduler struct {
	userService       user.Service
	host              string
	port              int
	url               string
	usersFreqInterval int
	limit             int
	logger            logger.Logger
}

func New(userService user.Service, host string, port int, url string, usersFreqInterval int, limit int, logger logger.Logger) *Scheduler {
	return &Scheduler{
		userService:       userService,
		host:              host,
		port:              port,
		url:               url,
		usersFreqInterval: usersFreqInterval,
		limit:             limit,
		logger:            logger,
	}
}

func (s *Scheduler) ScheduleUsers(ctx context.Context) {
	pingTimer := time.NewTicker(time.Duration(s.usersFreqInterval) * time.Minute)
	defer pingTimer.Stop()

	for {
		select {
		case <-pingTimer.C:
			lastRegisteredUser, err := s.userService.GetLastRegisteredUser(ctx)
			var lastRegisteredUserDate string
			if err == apperrors.ErrUserNotFound {
				lastRegisteredUserDate = "2000-04-23 10:03:32.670268"
			} else {
				lastRegisteredUserDate = lastRegisteredUser.CreatedAt.String()[0:26]
				s.logger.Info(fmt.Sprintf("Last registered user register date: %s", lastRegisteredUserDate))
			}

			users, err := s.getNewUsers(lastRegisteredUserDate)
			if err != nil {
				s.logger.Error(fmt.Sprintf("Error getting users in scheduler: %s", err.Error()))
				continue
			}

			if users != nil {
				s.logger.Info(fmt.Sprintf("%v", users))
				for _, u := range users {
					err := s.userService.Save(ctx, &u)
					if err != nil {
						s.logger.Warning(fmt.Sprintf("Error while saving new user [got in scheduler]: %s", err))
					}
				}
			} else {
				s.logger.Info("No more users")
			}
		}
	}
}

func (s *Scheduler) getNewUsers(lastRegisteredUserDate string) ([]user.User, error) {
	url := "http://" + s.host + ":" + strconv.Itoa(s.port) + s.url + "?date=" + lastRegisteredUserDate + "&limit=" + strconv.Itoa(s.limit)
	url = strings.Replace(url, " ", "%20", -1)

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
