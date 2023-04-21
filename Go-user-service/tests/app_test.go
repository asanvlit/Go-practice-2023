package tests

import (
	"Golang-practice-2023/internal/domain/user"
	"Golang-practice-2023/pkg/logger"
	"Golang-practice-2023/tests/data"
	"Golang-practice-2023/tests/data/provider"
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	zeroLogLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	myLogger, _ := logger.New(os.Getenv("LOG_LEVEL"), &zeroLogLogger)

	db, err := NewDb()
	if err != nil {
		myLogger.Fatal(err.Error())
	}

	userRepository, err := NewUserRepository(db, myLogger)
	if err != nil {
		myLogger.Fatal(err.Error())
	}

	userDataProvider, _ := NewUserDataProvider()

	userService, err := NewUserService(userRepository, myLogger)

	router := mux.NewRouter()
	port := os.Getenv("PORT")
	err = http.ListenAndServe(":"+port, router)

	t.Run("repository tests", func(t *testing.T) {
		RunRepositoryTests(userRepository, t)
	})
	t.Run("service tests", func(t *testing.T) {
		RunServiceTests(userService, userDataProvider, t)
	})

	t.Cleanup(func() {

	})
}

func RunRepositoryTests(repo user.Repository, t *testing.T) {
	t.Run("create-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := data.TestUser1()
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)

		query := "SELECT id, email, passwordhash FROM account WHERE id=$1"
		var returnedUser user.User
		err = repo.GetDbInstance().GetContext(ctx, &returnedUser, query, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, testUser.Email, returnedUser.Email)
		assert.Equal(t, testUser.Passwordhash, returnedUser.Passwordhash)

		repo.Delete(ctx, testUser.ID) // todo refactor
	})
	t.Run("get-by-id", func(t *testing.T) {
		ctx := context.Background()

		testUser := data.TestUser1()
		err := repo.Create(ctx, testUser)

		returnedUser, err := repo.GetById(ctx, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, testUser.Email, returnedUser.Email)
		assert.Equal(t, testUser.Passwordhash, returnedUser.Passwordhash)

		repo.Delete(ctx, testUser.ID) // todo refactor
	})
	t.Run("get-by-id-with-invalid-id", func(t *testing.T) {
		ctx := context.Background()
		invalidId, _ := uuid.Parse("da986a08-1aba-406c-9d56-bcf653e4d865")

		returnedUser, _ := repo.GetById(ctx, invalidId)

		require.Nil(t, returnedUser)
	})
	t.Run("update-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := data.TestUser1()
		_ = repo.Create(ctx, testUser)

		testUser2 := data.TestUser2()
		testUser.Email = testUser2.Email
		testUser.Passwordhash = testUser2.Passwordhash

		repo.Update(ctx, testUser)

		updatedTestUser, _ := repo.GetById(ctx, testUser.ID)
		assert.Equal(t, testUser2.Email, updatedTestUser.Email)
		assert.Equal(t, testUser2.Passwordhash, updatedTestUser.Passwordhash)

		repo.Delete(ctx, testUser.ID)
	})
	t.Run("delete-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := data.TestUser1WithId()
		_ = repo.Create(ctx, testUser)

		_ = repo.Delete(ctx, testUser.ID)

		returnedUser, _ := repo.GetById(ctx, testUser.ID)

		require.Nil(t, returnedUser)
	})
}

func RunServiceTests(service user.Service, provider *provider.UserDataProvider, t *testing.T) {
	t.Run("create-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		err := service.Create(ctx, testUser)
		require.NoError(t, err)

		createdUser, _ := service.GetById(ctx, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, testUser.Email, createdUser.Email)
		assert.Equal(t, testUser.Passwordhash, createdUser.Passwordhash)

		service.Delete(ctx, testUser.ID)
	})
	t.Run("create-user-with-invalid-email", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		testUser.Email = "testusergmail"
		err := service.Create(ctx, testUser)
		require.Error(t, err)
	})
	t.Run("create-user-with-invalid-password", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		testUser.Passwordhash = "xx"
		err := service.Create(ctx, testUser)
		require.Error(t, err)
	})
	t.Run("get-user-by-id", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		err := service.Create(ctx, testUser)

		returnedUser, err := service.GetById(ctx, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, testUser.Email, returnedUser.Email)
		assert.Equal(t, testUser.Passwordhash, returnedUser.Passwordhash)

		service.Delete(ctx, testUser.ID)
	})
	t.Run("get-user-by-invalid-id", func(t *testing.T) {
		ctx := context.Background()

		_, err := service.GetById(ctx, uuid.New())

		require.Error(t, err)
	})
	t.Run("update-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		_ = service.Create(ctx, testUser)

		testUser2 := provider.GenerateUserData(false, false)
		testUser.Email = testUser2.Email
		testUser.Passwordhash = testUser2.Passwordhash

		service.Update(ctx, testUser)

		updatedTestUser, _ := service.GetById(ctx, testUser.ID)
		assert.Equal(t, testUser2.Email, updatedTestUser.Email)

		service.Delete(ctx, testUser.ID)
	})
	t.Run("update-user-with-invalid-email", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		_ = service.Create(ctx, testUser)

		testUser2 := provider.GenerateUserData(false, false)
		testUser2.Email = "testusergmail"
		testUser.Email = testUser2.Email
		testUser.Passwordhash = testUser2.Passwordhash

		service.Update(ctx, testUser)

		updatedTestUser, _ := service.GetById(ctx, testUser.ID)
		assert.NotEqual(t, testUser2.Email, updatedTestUser.Email)

		service.Delete(ctx, testUser.ID)
	})
	t.Run("update-user-with-invalid-password", func(t *testing.T) {
		ctx := context.Background()

		testUser := provider.GenerateUserData(false, false)
		_ = service.Create(ctx, testUser)

		testUser2 := provider.GenerateUserData(false, false)
		testUser2.Passwordhash = "xx"
		testUser.Email = testUser2.Email
		testUser.Passwordhash = testUser2.Passwordhash

		service.Update(ctx, testUser)

		updatedTestUser, _ := service.GetById(ctx, testUser.ID)
		assert.NotEqual(t, testUser2.Email, updatedTestUser.Email)

		service.Delete(ctx, testUser.ID)
	})
	t.Run("delete-user", func(t *testing.T) {
		ctx := context.Background()

		testUser := data.TestUser1WithId()
		_ = service.Create(ctx, testUser)

		_ = service.Delete(ctx, testUser.ID)

		returnedUser, _ := service.GetById(ctx, testUser.ID)

		require.Nil(t, returnedUser)
	})
	t.Run("delete-user-with-invalid-id", func(t *testing.T) {
		ctx := context.Background()

		err := service.Delete(ctx, uuid.New())

		require.Error(t, err)
	})
}

const contentType = "application/json"
