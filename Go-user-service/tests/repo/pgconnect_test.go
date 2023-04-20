package repo

import (
	"Golang-practice-2023/pkg/pgconnect"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectDatabase(t *testing.T) {
	envPath, envErr := filepath.Abs("../../configs/test.env")
	if envErr != nil {
		t.Log("Can't get environment file") // todo t.log
	}

	err := godotenv.Load(envPath)
	if err != nil {
		t.Log("Error loading .env file")
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

	require.NotNil(t, db)
	require.NoError(t, err)
}

func TestConnectDatabaseFailOnWrongCredentials(t *testing.T) {
	config := pgconnect.ConnectionConfigData{
		Username:     "wronguser",
		Password:     "wrongpass",
		DatabaseName: "testdb",
		Port:         5432,
		Host:         "localhost",
	}

	db, err := pgconnect.ConnectDatabase(config)

	require.Error(t, err)
	require.Nil(t, db)
}

func TestConnect(t *testing.T) {
	t.Run("positive", func(t *testing.T) { TestConnectDatabase(t) })
}
