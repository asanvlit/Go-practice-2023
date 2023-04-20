package repo

import (
	"Golang-practice-2023/internal/domain/user"
	"context"
	_ "database/sql"
	_ "errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestComponent struct {
	repository user.Repository
}

func (tc *TestComponent) TestRepositoryCreate(ctx context.Context, repository user.Repository, t *testing.T) {
	testUser := testUser1()
	err := tc.repository.Create(ctx, testUser)
	require.NoError(t, err)

	query := "SELECT id, email, passwordhash FROM account WHERE id=$1"
	var returnedUser user.User
	err = tc.repository.GetDbInstance().GetContext(ctx, &returnedUser, query, testUser.ID)

	require.NoError(t, err)
	assert.Equal(t, testUser.Email, returnedUser.Email)
	assert.Equal(t, testUser.Passwordhash, returnedUser.Passwordhash)
}

func (tc *TestComponent) TestRepositoryGetById(ctx context.Context, t *testing.T) {
	testUser := testUser1()
	err := tc.repository.Create(ctx, testUser)

	query := "SELECT id, email, passwordhash FROM account WHERE id=$1"
	var returnedUser user.User
	err = tc.repository.GetDbInstance().GetContext(ctx, &returnedUser, query, testUser.ID)

	require.NoError(t, err)
	assert.Equal(t, testUser.Email, returnedUser.Email)
	assert.Equal(t, testUser.Passwordhash, returnedUser.Passwordhash)
}

func (tc *TestComponent) TestRepositoryGetByIdFailOnInvalidId(ctx context.Context, t *testing.T) {
	invalidId := "e658615b-66fe-49b5-85b6-519cda99495"

	query := "SELECT id, email, passwordhash FROM account WHERE id=$1"
	var returnedUser user.User
	err := tc.repository.GetDbInstance().GetContext(ctx, &returnedUser, query, invalidId)

	require.Error(t, err)
}

func (tc *TestComponent) TestRepositoryUpdate(ctx context.Context, t *testing.T) {
	testUser := testUser1()
	err := tc.repository.Create(ctx, testUser)
	testUser2 := testUser2()

	query := "UPDATE account SET email=$1, passwordhash=$2 WHERE id=$3"
	var returnedUser user.User
	err = tc.repository.GetDbInstance().GetContext(ctx, &returnedUser, query, testUser2.Email, testUser2.Passwordhash, testUser1().ID)

	require.NoError(t, err)
	assert.Equal(t, testUser2.Email, returnedUser.Email)
	assert.Equal(t, testUser2.Passwordhash, returnedUser.Passwordhash)
}

func (tc *TestComponent) TestRepositoryDelete(ctx context.Context, t *testing.T) {
	testUser := testUser1WithId()
	err := tc.repository.Create(ctx, testUser)

	query := "DELETE FROM account WHERE id=$3"
	var returnedUser user.User
	err = tc.repository.GetDbInstance().GetContext(ctx, &returnedUser, query, testUser1().ID)

	query2 := "SELECT id, email, passwordhash FROM account WHERE id=$1"
	var returnedUser2 user.User
	err = tc.repository.GetDbInstance().GetContext(ctx, &returnedUser2, query2, testUser1().ID)

	require.NoError(t, err)
	assert.Equal(t, returnedUser2, nil)
}

func testUser1() *user.User {
	return &user.User{
		ID:           uuid.Nil,
		Email:        "test1@gmail.com",
		Passwordhash: "11e176685d66625f153e7de7d547d77bf5797c7cca3aca713aa881b172da6613",
	}
}

func testUser1WithId() *user.User {
	id, _ := uuid.Parse("1f09ea98-0b87-4635-abf6-4cea3e8ea402")
	return &user.User{
		ID:           id,
		Email:        "test1@gmail.com",
		Passwordhash: "11e176685d66625f153e7de7d547d77bf5797c7cca3aca713aa881b172da6613",
	}
}

func testUser2() *user.User {
	return &user.User{
		ID:           uuid.Nil,
		Email:        "test2@gmail.com",
		Passwordhash: "369650404986c4bde35532defa1d74e7700b080349e811e24a352ea70150874c",
	}
}

func testUser2WithId() *user.User {
	id, _ := uuid.Parse("4a84c735-e026-43f4-af95-f9926b269d9f")
	return &user.User{
		ID:           id,
		Email:        "test2@gmail.com",
		Passwordhash: "369650404986c4bde35532defa1d74e7700b080349e811e24a352ea70150874c",
	}
}
