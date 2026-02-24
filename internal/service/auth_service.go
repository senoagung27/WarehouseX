package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/config"
	"github.com/senoagung27/warehousex/internal/domain/repository"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ AuthServiceInterface = (*AuthService)(nil)

type AuthService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
	log      *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtCfg config.JWTConfig, log *zap.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
		log:      log,
	}
}

func (s *AuthService) Register(input dto.RegisterInput) (*dto.AuthResponse, error) {
	existing, _ := s.userRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         input.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	s.log.Info("User registered", zap.String("email", user.Email), zap.String("role", user.Role))

	return &dto.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(input dto.LoginInput) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	s.log.Info("User logged in", zap.String("email", user.Email))

	return &dto.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.jwtCfg.ExpirationHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}
