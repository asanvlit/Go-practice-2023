package service

import (
	"Golang-practice-2023/internal/domain/logger"
	"Golang-practice-2023/internal/domain/user"
	"Golang-practice-2023/pkg/pubsub/nats/pub"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"regexp"
)

type Service struct {
	repository user.Repository
	logger     logger.Logger
	nats       *pub.NatsPublisher
}

func New(repository user.Repository, logger logger.Logger, nats *pub.NatsPublisher) *Service {
	return &Service{repository: repository, logger: logger, nats: nats}
}

func (service *Service) Create(ctx context.Context, user *user.User) error {
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if err := validatePassword(user.Passwordhash); err != nil {
		return err
	}

	hashedPassword := hashPassword(user.Passwordhash)
	user.Passwordhash = hashedPassword

	err := service.repository.Create(ctx, user)
	if err != nil {
		return err
	}

	//createdUserData := []byte(fmt.Sprintf("%s %s", user.ID, user.Email)) todo remove
	createdUserData, err := json.Marshal(user)
	if err != nil {
		service.logger.Warning("Failed to marshal user")
	}

	err = service.nats.Publish("NewUser", createdUserData)
	if err != nil {
		service.logger.Warning("Failed to push user data to nats")
	}

	return err
}

func (service *Service) GetById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return service.repository.GetById(ctx, id)
}

func (service *Service) Update(ctx context.Context, user *user.User) error {
	_, err := service.GetById(ctx, user.ID)
	if err != nil {
		return err
	}

	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if err := validatePassword(user.Passwordhash); err != nil {
		return err
	}

	hashedPassword := hashPassword(user.Passwordhash)
	user.Passwordhash = hashedPassword

	return service.repository.Update(ctx, user)
}

func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := service.GetById(ctx, id)
	if err != nil {
		return err
	}
	return service.repository.Delete(ctx, id)
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is empty")
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("password is empty")
	}

	passwordRegex := regexp.MustCompile("^[a-zA-Z0-9@#$%!]{8,60}$") // fixme reg exp
	if !passwordRegex.MatchString(password) {
		return errors.New("invalid password format")
	}

	return nil
}

func hashPassword(password string) string {
	hashed := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashed[:])
}
