package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const (
	RefreshTokenTTL = 7 * 24 * time.Hour
	ResetTokenTTL   = 1 * time.Hour
)

type UserService struct {
	userRepo    *UserRepo
	sessionRepo *SessionRepo
}

func NewUserService(userRepo *UserRepo, sessionRepo *SessionRepo) *UserService {
	return &UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *UserService) RegisterUser(user *User, password string) error {
	if check, _ := s.userRepo.GetUserByUsername(user.Username); check != nil {
		return ErrUserAlreadyExists
	}
	if check, _ := s.userRepo.GetUserByEmail(user.Email); check != nil {
		return ErrUserAlreadyExists
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	return s.userRepo.CreateUser(user)
}

// Login verifies credentials and issues a new access token + refresh session.
func (s *UserService) Login(username, password, ipAddress string) (user *User, accessToken string, refreshToken string, err error) {
	user, err = s.userRepo.GetUserByUsername(username)
	if err != nil || user == nil || !user.CheckPassword(password) {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, err = GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err = s.issueSession(user.ID, RefreshToken, ipAddress, RefreshTokenTTL)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Refresh(refreshToken, ipAddress string) (accessToken string, newRefreshToken string, err error) {
	session, err := s.sessionRepo.GetSessionByToken(refreshToken, RefreshToken)
	if err != nil {
		return "", "", err
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.DeleteSession(session.ID.String())
		return "", "", ErrSessionExpired
	}

	if err := s.sessionRepo.DeleteSession(session.ID.String()); err != nil {
		return "", "", err
	}

	accessToken, err = GenerateAccessToken(session.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.issueSession(session.UserID, RefreshToken, ipAddress, RefreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *UserService) Logout(refreshToken string) error {
	session, err := s.sessionRepo.GetSessionByToken(refreshToken, RefreshToken)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil
		}
		return err
	}
	return s.sessionRepo.DeleteSession(session.ID.String())
}

func (s *UserService) ForgotPassword(email, ipAddress string) (resetToken string, err error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", nil
	}

	return s.issueSession(user.ID, ResetPasswordToken, ipAddress, ResetTokenTTL)
}

func (s *UserService) ResetPassword(resetToken, newPassword string) error {
	session, err := s.sessionRepo.GetSessionByToken(resetToken, ResetPasswordToken)
	if err != nil {
		return err
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.DeleteSession(session.ID.String())
		return ErrSessionExpired
	}

	user, err := s.userRepo.GetUserByID(session.UserID.String())
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}
	user.UpdatedAt = time.Now()
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	if err := s.sessionRepo.DeleteSession(session.ID.String()); err != nil {
		return err
	}
	return s.sessionRepo.DeleteSessionsByUser(user.ID.String(), RefreshToken)
}

func (s *UserService) issueSession(userID uuid.UUID, tokenType TokenType, ipAddress string, ttl time.Duration) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}

	now := time.Now()
	session := &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		TokenType: tokenType,
		IpAddress: ipAddress,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}
