package service

import (
	"time"

	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
	"github.com/IndalAwalaikal/warung-pos/backend/utils"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
    Register(u *model.User) error
    // Login returns token and the authenticated user (or nil if invalid credentials)
    Login(email, password string) (string, *model.User, error)
}

type authService struct{
    userRepo repository.UserRepository
}

func NewAuthService(ur repository.UserRepository) AuthService {
    return &authService{userRepo: ur}
}

func (s *authService) Register(u *model.User) error {
    hashed, err := utils.HashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hashed
    if u.Role == "" {
        u.Role = "user"
    }
    return s.userRepo.Create(u)
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, err
    }
    if user == nil {
        return "", nil, nil
    }
    if !utils.CheckPassword(user.Password, password) {
        return "", nil, nil
    }

    // Create token
    claims := config.Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JwtExpiry())),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(config.JwtSecret())
    if err != nil {
        return "", nil, err
    }
    return signed, user, nil
}
