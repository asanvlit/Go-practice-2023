package tests

import (
	"Golang-practice-2023/internal/domain/logger"
	"Golang-practice-2023/internal/domain/user"
	"Golang-practice-2023/internal/user/repository"
	"Golang-practice-2023/internal/user/service"
	"Golang-practice-2023/pkg/migration"
	"Golang-practice-2023/pkg/pgconnect"
	"Golang-practice-2023/tests/data/provider"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func NewDb() (*sqlx.DB, error) {
	err := godotenv.Load("C:\\Go\\GoProjects\\Golang-tech-practice-2023\\Golang-practice-2023\\configs\\test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("Failed to get Postgresql post")
	}

	db, err := pgconnect.ConnectDatabase(pgconnect.ConnectionConfigData{
		Username:     os.Getenv("POSTGRES_USERNAME"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("POSTGRES_DATABASE"),
		Port:         pgPort,
		Host:         os.Getenv("POSTGRES_HOST"),
	})
	_ = db

	migrationData, err := migration.New(db.DB, "file:../schemas")
	err = migrationData.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	} else {
		log.Print(err)
	}

	return db, nil
}

func NewUserRepository(db *sqlx.DB, logger logger.Logger) (*repository.Repository, error) {
	return repository.New(db, logger), nil
}

func NewUserService(repository user.Repository, logger logger.Logger) (*service.Service, error) {
	return service.New(repository, logger), nil
}

func NewUserDataProvider() (*provider.UserDataProvider, error) {
	return provider.New(), nil
}
