package service

import (
	"context"
	"crypto/rand"
	"docs_storage/internal/models"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
    Create(ctx context.Context, u *models.User) error
    GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type sessionRepository interface {
    Create(ctx context.Context, s *models.Session) error
    GetByToken(ctx context.Context, token string) (*models.Session, error)
    Delete(ctx context.Context, token string) error
}

type AuthService struct {
    users    userRepository
    sessions sessionRepository
    adminTok string
}

func NewAuthService(users userRepository, sessions sessionRepository, adminToken string) *AuthService {
    return &AuthService{users: users, sessions: sessions, adminTok: adminToken}
}

func (s *AuthService) Register(ctx context.Context, token, login, pswd string) error {
    if token != s.adminTok {
        return errors.New("unauthorized")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    u := &models.User{
        Login:    login,
        Password: string(hash),
    }
    return s.users.Create(ctx, u)
}

func (s *AuthService) Auth(ctx context.Context, login, pswd string) (string, error) {
    u, err := s.users.GetByLogin(ctx, login)
    if err != nil {
        return "", err
    }
    if u == nil {
        return "", errors.New("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pswd)); err != nil {
        return "", errors.New("invalid credentials")
    }

    buf := make([]byte, 16)
    rand.Read(buf)
    token := hex.EncodeToString(buf)

    sess := &models.Session{
        Token:  token,
        UserID: u.ID,
        Login:  u.Login,
    }

    if err := s.sessions.Create(ctx, sess); err != nil {
        return "", err
    }

    return token, nil
}


func (s *AuthService) Logout(ctx context.Context, token string) error {
    return s.sessions.Delete(ctx, token)
}
