package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hrustalq/go_backend/proto/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer
	DB        *gorm.DB
	Redis     *redis.Client
	JWTSecret []byte
}

func (s *AuthService) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &auth.SignUpResponse{Message: "User created successfully"}, nil
}

func (s *AuthService) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	var user User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &auth.SignInResponse{Token: tokenString}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	return &auth.ValidateTokenResponse{Valid: true}, nil
}
