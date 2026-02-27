package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error
	IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, id uuid.UUID) error
	LockAccount(ctx context.Context, id uuid.UUID, until time.Time) error

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new userRepository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, role, status,
		       last_login_at, last_login_ip, failed_attempts, locked_until,
		       email_verified, metadata, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, role, status,
		       last_login_at, last_login_ip, failed_attempts, locked_until,
		       email_verified, metadata, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}
	return &user, nil
}

// FindAll retrieves all users with optional filtering
func (r *userRepository) FindAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.Role != nil {
		where += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, *filter.Role)
		argIdx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (email ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx+1)
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("userRepository.FindAll count: %w", err)
	}

	// Data query
	filter.Normalize()
	dataQuery := fmt.Sprintf(`
		SELECT id, email, password_hash, full_name, avatar_url, role, status,
		       last_login_at, last_login_ip, failed_attempts, locked_until,
		       email_verified, metadata, created_at, updated_at, deleted_at
		FROM users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var users []*domain.User
	if err := r.db.SelectContext(ctx, &users, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("userRepository.FindAll: %w", err)
	}

	return users, total, nil
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, avatar_url, role, status, email_verified, metadata)
		VALUES (:id, :email, :password_hash, :full_name, :avatar_url, :role, :status, :email_verified, :metadata)
		RETURNING created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.CreatedAt, &user.UpdatedAt); err != nil {
			return fmt.Errorf("userRepository.Create scan: %w", err)
		}
	}
	return nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET full_name = :full_name, avatar_url = :avatar_url, role = :role,
		    status = :status, email_verified = :email_verified, metadata = :metadata,
		    updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("userRepository.Update: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.UpdatedAt); err != nil {
			return fmt.Errorf("userRepository.Update scan: %w", err)
		}
	}
	return nil
}

// Delete soft-deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp and IP
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
	query := `UPDATE users SET last_login_at = NOW(), last_login_ip = $1, failed_attempts = 0 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, ip, id)
	if err != nil {
		return fmt.Errorf("userRepository.UpdateLastLogin: %w", err)
	}
	return nil
}

// IncrementFailedAttempts increments the failed login attempts counter
func (r *userRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET failed_attempts = failed_attempts + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepository.IncrementFailedAttempts: %w", err)
	}
	return nil
}

// ResetFailedAttempts resets the failed login attempts counter
func (r *userRepository) ResetFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET failed_attempts = 0, locked_until = NULL WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepository.ResetFailedAttempts: %w", err)
	}
	return nil
}

// LockAccount locks a user account until the specified time
func (r *userRepository) LockAccount(ctx context.Context, id uuid.UUID, until time.Time) error {
	query := `UPDATE users SET locked_until = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, until, id)
	if err != nil {
		return fmt.Errorf("userRepository.LockAccount: %w", err)
	}
	return nil
}

// CreateRefreshToken stores a new refresh token
func (r *userRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, ip_address, user_agent)
		VALUES (:id, :user_id, :token_hash, :expires_at, :ip_address, :user_agent)
		RETURNING created_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("userRepository.CreateRefreshToken: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&token.CreatedAt); err != nil {
			return fmt.Errorf("userRepository.CreateRefreshToken scan: %w", err)
		}
	}
	return nil
}

// FindRefreshToken retrieves a refresh token by its hash
func (r *userRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, ip_address, user_agent, created_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var token domain.RefreshToken
	if err := r.db.GetContext(ctx, &token, query, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindRefreshToken: %w", err)
	}
	return &token, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = NOW() WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("userRepository.RevokeRefreshToken: %w", err)
	}
	return nil
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (r *userRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = NOW() WHERE user_id = $1 AND is_revoked = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("userRepository.RevokeAllUserRefreshTokens: %w", err)
	}
	return nil
}

// DeleteExpiredRefreshTokens removes expired refresh tokens from the database
func (r *userRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("userRepository.DeleteExpiredRefreshTokens: %w", err)
	}
	return nil
}
