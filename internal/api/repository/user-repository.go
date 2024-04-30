package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-fitness/external/db"
	"go-fitness/internal/api/types"
)

type UserRepository struct {
	db db.SqlInterface
}

type UserRepositoryInterface interface {
	GetUserByUUID(ctx context.Context, uuid string) (types.User, error)
	GetRoleByUserID(ctx context.Context, userID int64) (types.Role, error)
}

func NewUserRepository(db db.SqlInterface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByUUID(ctx context.Context, uuid string) (types.User, error) {
	const op = "repository.user.FindUserByUUID"

	const query = "SELECT id,uuid,name,email FROM users WHERE uuid = ? AND active = 1 AND deleted_at IS NULL AND email_verified_at IS NOT NULL"

	row := r.db.GetExecer().QueryRowContext(ctx, query, uuid)

	user := types.User{}

	err := row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, err)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) GetRoleByUserID(ctx context.Context, userID int64) (types.Role, error) {
	const op = "repository.user.GetRoleByUserID"

	const query = "SELECT r.id, r.name FROM roles r INNER JOIN model_has_roles ur ON r.id = ur.role_id WHERE ur.model_id = ?"

	row := r.db.GetExecer().QueryRowContext(ctx, query, userID)

	role := types.Role{}

	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return role, fmt.Errorf("%s: %w", op, err)
		}

		return role, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
