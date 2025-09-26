package seeder

import "faizalmaulana/lsp/models/repo"

func RunAll(users repo.UsersRepo, profiles repo.ProfilesRepo) error {
	if err := SeedAdmin(users, profiles); err != nil {
		return err
	}
	return nil
}
