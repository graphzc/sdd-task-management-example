package user

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
	"github.com/graphzc/sdd-task-management-example/internal/infrastructure/auth"
	"github.com/graphzc/sdd-task-management-example/internal/repositories/user"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/graphzc/sdd-task-management-example/internal/utils/timeutil"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, in *UserRegisterInput) error
	Login(ctx context.Context, in *UserLoginInput) (string, error)
}

type service struct {
	config *config.Config
	repo   user.Repository
}

// @WireSet("Service")
func NewService(
	config *config.Config,
	repo user.Repository,
) Service {
	return &service{
		config: config,
		repo:   repo,
	}
}

func (s *service) Register(ctx context.Context, in *UserRegisterInput) error {
	// Check if user with the same email already exists
	existingUser, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to check existing user by email")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to check existing user",
		)
	}

	if existingUser != nil {
		log.Warn().
			Str("email", in.Email).
			Msg("User with the same email already exists")

		return servererr.NewError(
			servererr.ErrorCodeConflict,
			"User with the same email already exists",
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user
	newUser := &entities.User{
		ID:        uuid.NewString(),
		Email:     in.Email,
		Password:  string(hashedPassword),
		Name:      in.Name,
		CreatedAt: timeutil.BangkokNow(),
		UpdatedAt: timeutil.BangkokNow(),
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create new user")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to create new user",
		)
	}

	return nil
}

func (s *service) Login(ctx context.Context, in *UserLoginInput) (string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to find user by email")

		return "", servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to find user",
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		log.Warn().
			Str("email", in.Email).
			Str("password", in.Password).
			Msg("Invalid password")

		return "", servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"Invalid email or password",
		)
	}

	accessTokenExpiration, err := time.ParseDuration(s.config.JWT.AccessTokenExpiration)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse access token expiration duration")

		return "", servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to create access token",
		)
	}

	expiredAt := timeutil.BangkokNow().Add(accessTokenExpiration)

	// Generate JWT token
	token, err := generateJWTToken(user, expiredAt, s.config.JWT.AccessTokenSecret)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to generate JWT token")

		return "", servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to generate JWT token",
		)
	}

	return token, nil
}

func generateJWTToken(user *entities.User, expiredAt time.Time, secret string) (string, error) {
	claims := auth.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "task-management",
			Subject:   user.ID,
			Audience:  []string{"task-management-users"},
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			NotBefore: jwt.NewNumericDate(timeutil.BangkokNow()),
			IssuedAt:  jwt.NewNumericDate(timeutil.BangkokNow()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
