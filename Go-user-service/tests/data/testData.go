package data

import (
	"Golang-practice-2023/internal/domain/user"
	"github.com/google/uuid"
)

func TestUser1() *user.User {
	return &user.User{
		ID:           uuid.Nil,
		Email:        "test11@gmail.com",
		Passwordhash: "11e176685d66625f153e7de7d547d77bf5797c7cca3aca713aa881b172da6613",
	}
}

func TestUser1WithId() *user.User {
	id, _ := uuid.Parse("1f09ea98-0b87-4635-abf6-4cea3e8ea402")
	return &user.User{
		ID:           id,
		Email:        "test11@gmail.com",
		Passwordhash: "11e176685d66625f153e7de7d547d77bf5797c7cca3aca713aa881b172da6613",
	}
}

func TestUser2() *user.User {
	return &user.User{
		ID:           uuid.Nil,
		Email:        "test22@gmail.com",
		Passwordhash: "369650404986c4bde35532defa1d74e7700b080349e811e24a352ea70150874c",
	}
}
