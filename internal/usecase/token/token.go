package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
	"tender_management/internal/entity"
	"time"
)

var (
	AccessSecretKey  string
	RefreshSecretKey string
	ExpiredAccess    int
	ExpiredRefresh   int
)

func ValidateToken(tokenstr string) (bool, error) {
	_, err := ExtractClaims(tokenstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractClaims(tokenstr string) (jwt.MapClaims, error) {
	tokenstr = strings.TrimPrefix(tokenstr, "\"")
	tokenstr = strings.TrimSuffix(tokenstr, "\"")
	token, err := jwt.ParseWithClaims(tokenstr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(AccessSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenstr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("failed to parse token claims")
	}
	return claims, nil
}
func GenerateAccessToken(in entity.User) (string, error) {
	claims := Claims{
		Id:       in.ID,
		Username: in.Username,
		Role:     in.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpiredAccess)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(os.Getenv(AccessSecretKey)))

	return str, err
}

func GenerateRefreshToken(in entity.User) (string, error) {
	claims := Claims{
		Id:       in.ID,
		Username: in.Username,
		Role:     in.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpiredRefresh)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(os.Getenv(RefreshSecretKey)))

	return str, err
}

func GetExpires() int {
	return ExpiredAccess
}
