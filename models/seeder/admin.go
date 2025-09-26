package seeder

import (
	"time"

	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(users repo.UsersRepo) error {
	_, err := users.GetByEmail("mail@paisaltanjung.my.id")
	if err == nil {
		return nil
	}

	pwd := "admin123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &entity.Users{
		IdUser:    helper.Uuid(),
		Email:     "mail@paisaltanjung.my.id",
		Password:  string(hashed),
		Role:      "admin",
		Timestamp: time.Now(),
	}

	if err := users.Create(u); err != nil {
		return err
	}
	return nil
}
