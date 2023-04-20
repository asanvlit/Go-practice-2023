package repo

import (
	"Golang-practice-2023/pkg/migration"
	"Golang-practice-2023/pkg/pgconnect"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestAccountsMigrateUp(t *testing.T) {
	envPath, envErr := filepath.Abs("../../configs/test.env")
	if envErr != nil {
		log.Print(envErr)
		t.Log("Can't get environment file")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		t.Log("Error loading .env file")
	}

	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		t.Log("Failed to get Postgresql port")
	}
	db, err := pgconnect.ConnectDatabase(pgconnect.ConnectionConfigData{
		Username:     os.Getenv("POSTGRES_USERNAME"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("POSTGRES_DATABASE"),
		Port:         pgPort,
		Host:         os.Getenv("POSTGRES_HOST"),
	})

	var hasChanged bool
	migrationData, err := migration.New(db.DB, "file://../../schemas")
	err = migrationData.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Log(err)
		require.NoError(t, err)
	} else {
		t.Log("No changes")
	}
	if err == nil {
		hasChanged = true
	}

	rows, err := db.Query("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'account')")
	require.NoError(t, err)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Log("Failed to close")
		}
	}(rows)

	var exists bool
	for rows.Next() {
		err = rows.Scan(&exists)
		require.NoError(t, err)
	}
	assert.True(t, exists)

	if hasChanged {
		err = migrationData.Down()
	}
}
