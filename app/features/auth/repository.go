package auth

import (
	"context"

	"northstar/app/features/auth/gen/authdb"
)

type authRepository struct {
	queries *authdb.Queries
}

func (r *authRepository) getUserByEmail(ctx context.Context, email string) (authdb.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *authRepository) createUser(ctx context.Context, params authdb.CreateUserParams) (authdb.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *authRepository) checkIfUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := r.queries.CheckIfUserExistsByUsername(ctx, username)
	return exists != 0, err
}
func (r *authRepository) checkIfUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.CheckIfUserExistsByEmail(ctx, email)
	return exists != 0, err
}
