package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"
)

type AuthUsecase interface {
	Register(ctx context.Context, email, password string, displayName *string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*entity.User, error)
}

type authUsecase struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (a *authUsecase) Register(ctx context.Context, email, password string, displayName *string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := a.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entity.NewUser(email, string(hashedPassword), displayName)

	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, errors.New("account is inactive")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Update last login
	err = a.userRepo.UpdateLastLogin(ctx, user.ID, time.Now())
	if err != nil {
		// Log error but don't fail login
	}

	// Generate JWT token
	token, err := a.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (a *authUsecase) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	user, err := a.ValidateToken(ctx, tokenString)
	if err != nil {
		return "", err
	}

	return a.generateToken(user)
}

func (a *authUsecase) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	user, err := a.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	return user, nil
}

func (a *authUsecase) generateToken(user *entity.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}
