package service

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type AuthService struct {
	users  *repository.UserRepository
	tokens *support.TokenManager
}

func NewAuthService(users *repository.UserRepository, tokens *support.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Register(input dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.users.FindByEmail(input.Email); err == nil {
		return nil, support.NewError(http.StatusConflict, "email_taken", "email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Username:     input.Username,
		Email:        strings.ToLower(input.Email),
		PasswordHash: s.tokens.HashPassword(input.Password),
		Role:         model.RoleUser,
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	return s.issue(user), nil
}

func (s *AuthService) Login(input dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		}
		return nil, err
	}

	if !s.tokens.CheckPassword(user.PasswordHash, input.Password) {
		return nil, support.NewError(http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
	}

	return s.issue(user), nil
}

func (s *AuthService) Logout(userID uint) error {
	if _, err := s.users.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusUnauthorized, "unauthorized", "user session not found")
		}
		return err
	}
	return nil
}

func (s *AuthService) Profile(userID uint) (*dto.UserSummary, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *AuthService) issue(user *model.User) *dto.AuthResponse {
	return &dto.AuthResponse{
		Token: s.tokens.Issue(*user),
		User:  toUserSummary(*user),
	}
}

func toUserSummary(user model.User) dto.UserSummary {
	return dto.UserSummary{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: model.RolePermissions(user.Role),
		Bio:         user.Bio,
		Avatar:      user.Avatar,
	}
}
