package jwt

import (
	"time"
	config "user-service/config/grpc_server"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

func NewTokenManager(config *config.Config) *TokenManager {
	return &TokenManager{
		SecretKey:     config.JWTSecret,
		AccessExpiry:  config.AccessExpiration,
		RefreshExpiry: config.RefreshExpiration,
		Issuer:        config.TokenIssuer,
	}
}

type Claims struct {
	jwt.StandardClaims
	UserID int64  `json:"user_id"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
}

func (m *TokenManager) GenerateAccessToken(userID int64, email, role string) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(m.AccessExpiry).Unix(),
			Issuer:    m.Issuer,
		},
		UserID: userID,
		Email:  email,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.SecretKey))
}

func (m *TokenManager) GenerateRefreshToken(userID int64) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(m.RefreshExpiry).Unix(),
			Issuer:    m.Issuer,
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.SecretKey))
}

func (m *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(m.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (m *TokenManager) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func (m *TokenManager) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Генерируем новые токены
	accessToken, err := m.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := m.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
