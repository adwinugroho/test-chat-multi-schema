package service

import (
	"context"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) generateJWTToken(userID, role string) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_role": role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":       time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		logger.LogError("Error while get user by email:" + err.Error())
		return nil, model.NewError(model.ErrorGeneral, "Internal Server Error")
	}

	if user == nil {
		logger.LogInfo("User not found")
		return nil, model.NewError(model.ErrorDataFound, "User Not Found")
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		logger.LogError("Error while get user by ID:" + err.Error())
		return nil, model.NewError(model.ErrorGeneral, "Internal Server Error")
	}

	if user == nil {
		logger.LogInfo("User not found")
		return nil, model.NewError(model.ErrorDataFound, "User Not Found")
	}

	return user, nil
}

func (s *userService) LoginUser(ctx context.Context, req model.LoginUserRequest) (*model.AuthenticationResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.LogError("Login failed, cannot find user: " + err.Error())
		return nil, model.NewError(model.ErrorUnauthorized, "Invalid email or password")
	}
	if user == nil {
		return nil, model.NewError(model.ErrorUnauthorized, "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return nil, model.NewError(model.ErrorUnauthorized, "Invalid email or password")
	}

	// Generate token
	token, err := s.generateJWTToken(user.UserID, user.Role)
	if err != nil {
		logger.LogError("Failed to generate token: " + err.Error())
		return nil, model.NewError(model.ErrorGeneral, "Internal server error")
	}

	return &model.AuthenticationResponse{
		ID:    user.UserID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}, nil
}
