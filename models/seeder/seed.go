package seeder

import "faizalmaulana/lsp/models/repo"

// RunAll executes all seeders. Currently only seeds admin user.
func RunAll(users repo.UsersRepo) error {
	if err := SeedAdmin(users); err != nil {
		return err
	}
	return nil
}
