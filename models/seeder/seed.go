package seeder

import "faizalmaulana/lsp/models/repo"

func RunAll(users repo.UsersRepo) error {
	if err := SeedAdmin(users); err != nil {
		return err
	}
	return nil
}
