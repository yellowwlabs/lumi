package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
	ForgotPasswordToken
	ResetPasswordToken
)

type User struct {
	ID            uuid.UUID
	DisplayName   string
	Username      string
	PasswordHash  string
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	TokenType TokenType
	IpAddress string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
