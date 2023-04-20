package provider

import (
	"Golang-practice-2023/internal/domain/user"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lucasjones/reggen"
)

type UserDataProvider struct {
}

func New() *UserDataProvider {
	return &UserDataProvider{}
}

func (provider *UserDataProvider) GenerateUserData(hasId bool, isHashedPassword bool) *user.User {
	var id uuid.UUID
	if hasId {
		id = generateId()
	} else {
		id = uuid.Nil
	}

	password := generatePassword()
	if isHashedPassword {
		password = hashPassword(password)
	}

	return &user.User{
		ID:           id,
		Email:        generateEmail(),
		Passwordhash: password,
	}
}

func (provider *UserDataProvider) GenerateUserList(count int, hasId bool, isHashedPassword bool) []*user.User {
	userList := make([]*user.User, count)
	for i := 0; i < count; i++ {
		userList[i] = provider.GenerateUserData(hasId, isHashedPassword)
	}
	return userList
}

func generateId() uuid.UUID {
	return uuid.New()
}

func generateEmail() string {
	email, err := reggen.Generate("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", 25)
	if err != nil {
		panic(err)
	}
	return email
}

func generatePassword() string {
	password, err := reggen.Generate("^[a-zA-Z0-9@#$%!]{8,60}$", 25)
	if err != nil {
		panic(err)
	}
	return password
}

func hashPassword(password string) string {
	hashed := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashed[:])
}
