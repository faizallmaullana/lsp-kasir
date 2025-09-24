package seeder

import (
	"time"

	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"

	"golang.org/x/crypto/bcrypt"
)

// SeedAdmin ensures an admin user exists. If a user with username 'admin' does not exist, it will be created.
func SeedAdmin(users repo.UsersRepo) error {
	_, err := users.GetByEmail("mail@paisaltanjung.my.id")
	if err == nil {
		// admin already exists
		return nil
	}

	pwd := "admin123" // default password; consider taking from env
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
