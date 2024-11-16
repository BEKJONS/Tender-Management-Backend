package token

import (
	"github.com/golang-jwt/jwt"
	"os"
	"tender_management/internal/entity"
	"time"
)

var (
	AccessSecretKey  string
	RefreshSecretKey string
	ExpiredAccess    int
	ExpiredRefresh   int
)

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