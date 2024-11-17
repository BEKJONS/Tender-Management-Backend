package token

import (
	"fmt"
	"strconv"
	"tender_management/config"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func ConfigToken(config *config.Config) error {

	exAcc, err := strconv.Atoi(config.EXPIRED_ACCESS)
	if err != nil {
		return err
	}

	refAcc, err := strconv.Atoi(config.EXPIRED_REFRESH)
	if err != nil {
		return err
	}

	AccessSecretKey = config.ACCESS_TOKEN
	RefreshSecretKey = config.ACCESS_TOKEN
	ExpiredAccess = exAcc
	ExpiredRefresh = refAcc

	return nil
}

func ValidateToken(tokenstr string) (bool, error) {
	_, err := ExtractClaims(tokenstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractClaims(tokenstr string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Token imzosi HMAC bo'lishi kerak
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.NewConfig().ACCESS_TOKEN), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenstr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	return claims, nil
}
