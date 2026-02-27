package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleEditor     UserRole = "editor"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// User represents an admin user
type User struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Email          string     `db:"email" json:"email"`
	PasswordHash   string     `db:"password_hash" json:"-"`
	FullName       string     `db:"full_name" json:"full_name"`
	AvatarURL      *string    `db:"avatar_url" json:"avatar_url"`
	Role           UserRole   `db:"role" json:"role"`
	Status         UserStatus `db:"status" json:"status"`
	LastLoginAt    *time.Time `db:"last_login_at" json:"last_login_at"`
	LastLoginIP    *string    `db:"last_login_ip" json:"last_login_ip"`
	FailedAttempts int        `db:"failed_attempts" json:"-"`
	LockedUntil    *time.Time `db:"locked_until" json:"-"`
	EmailVerified  bool       `db:"email_verified" json:"email_verified"`
	Metadata       JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"-"`
}

// IsActive returns true if the user is active and not locked
func (u *User) IsActive() bool {
	return u.Status == StatusActive && u.DeletedAt == nil
}

// IsLocked returns true if the user account is temporarily locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// HasRole returns true if the user has the specified role or higher
func (u *User) HasRole(role UserRole) bool {
	switch role {
	case RoleEditor:
		return u.Role == RoleEditor || u.Role == RoleAdmin || u.Role == RoleSuperAdmin
	case RoleAdmin:
		return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
	case RoleSuperAdmin:
		return u.Role == RoleSuperAdmin
	}
	return false
}

// RefreshToken represents a JWT refresh token stored in the database
type RefreshToken struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	TokenHash string     `db:"token_hash" json:"-"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	IsRevoked bool       `db:"is_revoked" json:"is_revoked"`
	IPAddress *string    `db:"ip_address" json:"ip_address"`
	UserAgent *string    `db:"user_agent" json:"user_agent"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	RevokedAt *time.Time `db:"revoked_at" json:"revoked_at"`
}

// IsValid returns true if the token is not expired and not revoked
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && time.Now().Before(rt.ExpiresAt)
}

// UserFilter holds filter parameters for user queries
type UserFilter struct {
	Role   *UserRole
	Status *UserStatus
	Search *string
	Pagination
}

// CreateUserInput holds data for creating a new user
type CreateUserInput struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8,max=72"`
	FullName string   `json:"full_name" validate:"required,min=2,max=255"`
	Role     UserRole `json:"role" validate:"required,oneof=super_admin admin editor"`
}

// UpdateUserInput holds data for updating a user
type UpdateUserInput struct {
	FullName  *string    `json:"full_name" validate:"omitempty,min=2,max=255"`
	AvatarURL *string    `json:"avatar_url" validate:"omitempty,url"`
	Role      *UserRole  `json:"role" validate:"omitempty,oneof=super_admin admin editor"`
	Status    *UserStatus `json:"status" validate:"omitempty,oneof=active inactive suspended"`
}

// ChangePasswordInput holds data for changing a user's password
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}

// LoginInput holds data for user login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthTokens holds the access and refresh tokens
type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// JWTClaims holds the JWT claims
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   UserRole  `json:"role"`
}
