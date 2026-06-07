package repository

import (
	"context"
	"time"

	"game-platform/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository репозиторий для работы с пользователями
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// GetBySID получает пользователя по SID
func (r *UserRepository) GetBySID(ctx context.Context, sid string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT sid, email, name, gender, department, position, photo_url, last_sync, created_at, updated_at
		FROM users WHERE sid = $1
	`, sid).Scan(
		&user.SID, &user.Email, &user.Name, &user.Gender, &user.Department,
		&user.Position, &user.PhotoURL, &user.LastSync, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, sid, email, name, gender, department, position, photoURL string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (sid, email, name, gender, department, position, photo_url, last_sync)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING sid, email, name, gender, department, position, photo_url, last_sync, created_at, updated_at
	`, sid, email, name, gender, department, position, photoURL, time.Now()).Scan(
		&user.SID, &user.Email, &user.Name, &user.Gender, &user.Department,
		&user.Position, &user.PhotoURL, &user.LastSync, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update обновляет пользователя
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET name = $2, gender = $3, department = $4, position = $5, photo_url = $6, last_sync = $7, updated_at = NOW()
		WHERE sid = $1
	`, user.SID, user.Name, user.Gender, user.Department, user.Position, user.PhotoURL, time.Now())
	return err
}

// GetUserGroups получает группы пользователя
func (r *UserRepository) GetUserGroups(ctx context.Context, sid string) ([]string, error) {
	rows, err := r.db.Query(ctx, "SELECT group_name FROM user_groups WHERE sid = $1", sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
