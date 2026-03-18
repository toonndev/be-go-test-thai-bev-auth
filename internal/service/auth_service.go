package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"be-go-test-thai-bev-auth/internal/dto"
	"be-go-test-thai-bev-auth/internal/model"
	"be-go-test-thai-bev-auth/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrPasswordMismatch  = errors.New("password and confirm_password do not match")
	ErrUsernameExists    = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	fmt.Print("req =>",req)
	if req.Password != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return ErrUsernameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashed),
	}
	return s.userRepo.Create(user)
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:    token,
		Username: user.Username,
	}, nil
}

func generateJWT(username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
