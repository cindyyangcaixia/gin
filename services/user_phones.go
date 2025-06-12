package services

import (
	"context"
	stderrors "errors"
	"scalper/errors"
	"scalper/models"
	"scalper/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Claims struct {
	Name         string `json:"name,omitempty"`
	SerialNumber string `json:"serial_number"`
	PhoneNumber  string `json:"phone_number"`
	AmID         string `json:am_id`
	jwt.RegisteredClaims
}

type UserPhoneService struct {
	repo   *repositories.UserPhoneRepository
	logger *zap.Logger
}

func NewUserPhoneService(repo *repositories.UserPhoneRepository,
	logger *zap.Logger) *UserPhoneService {
	return &UserPhoneService{repo: repo, logger: logger}
}

func (s *UserPhoneService) CreateUserPhone(ctx context.Context,
	userPhone *models.UserPhone) (*mongo.InsertOneResult, error) {

	// hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	hashedPassword := string(hashedBytes)
	// }
	return s.repo.InsertOne(ctx, userPhone)
}

func (s *UserPhoneService) GetUserPhone(ctx context.Context,
	phoneNumber string) (*models.UserPhone, error) {
	return s.repo.FindOne(ctx, phoneNumber)
}

func (s *UserPhoneService) ListUserPhones(ctx context.Context, phoneNumber string,
	serialNumber *string, page, limit int64) ([]*models.UserPhone, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.ListUserPhones(ctx, phoneNumber, serialNumber, page, limit)
}

func (s *UserPhoneService) Login(ctx context.Context, phoneNumber string, password string) (string, error) {
	userPhone, err := s.repo.FindOne(ctx, phoneNumber)
	s.logger.Info("users: %v", zap.Error(err))
	if err != nil {
		return "", errors.NewAppError(errors.ErrCodeUserNotExist, stderrors.New("invalid users"), "")
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(userPhone.PasswordHash), []byte(password)); err != nil {
	// 	return "", errors.New("invalid credentials")
	// }

	expirationTime := time.Now().Add(time.Duration(viper.GetInt64("JwtExpiration")) * time.Hour)
	claims := &Claims{
		PhoneNumber:  userPhone.PhoneNumber,
		SerialNumber: userPhone.SerialNumber,
		AmID:         userPhone.AmID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("JwtSecret")))
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return "", errors.NewAppError(errors.ErrCodeInvalidToken, stderrors.New("failed to generate token"), "")
	}
	return tokenString, nil
}
