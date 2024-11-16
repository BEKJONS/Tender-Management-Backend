package token

import (
	"github.com/golang-jwt/jwt"
	"strconv"
	"tender_management/config"
)

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func ConfigToken(config config.Config) error {

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
