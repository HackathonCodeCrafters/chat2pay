package jwt

import (
	"chat2pay/config/yaml"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware interface {
	GenerateToken(userID string, email string, role string) (*string, error)
	GenerateTokenWithMerchant(userID string, merchantID string, email string, role string) (*string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID     string `json:"user_id"`
	MerchantID string `json:"merchant_id,omitempty"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

type authMiddleware struct {
	cfg *yaml.Config
}

func NewAuthMiddleware(cfg *yaml.Config) AuthMiddleware {
	return &authMiddleware{
		cfg: cfg,
	}
}

func (m *authMiddleware) GenerateToken(userID string, email string, role string) (*string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.cfg.JWT.Key))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (m *authMiddleware) GenerateTokenWithMerchant(userID string, merchantID string, email string, role string) (*string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:     userID,
		MerchantID: merchantID,
		Email:      email,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.cfg.JWT.Key))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (m *authMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.cfg.JWT.Key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
