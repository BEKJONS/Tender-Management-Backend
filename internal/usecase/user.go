package usecase

import (
	"fmt"
	"log/slog"
	"tender_management/internal/entity"
	"tender_management/internal/usecase/help"
	"tender_management/internal/usecase/token"
)

type UserUseCase struct {
	repo UsersRepo
	log  *slog.Logger
}

func NewUserUseCase(repo UsersRepo, log *slog.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log,
	}
}

func (u *UserUseCase) AddUser(in entity.RegisterReq) (*entity.LogInRes, error) {
	u.log.Info("Add user started", "username", in.Username)
	hash, err := help.HashPassword(in.Password)
	if err != nil {
		u.log.Error("Error in hashing password", "error", err)
		return nil, err
	}

	in.Password = hash
	if in.Role == "" {
		in.Role = "contractor"
	}

	res, err := u.repo.CreateUser(entity.User{Username: in.Username, Password: in.Password, Role: in.Role, Email: in.Email})
	if err != nil {
		u.log.Error("Error in adding user", "error", err)
		return nil, err
	}
	accessToken, err := token.GenerateAccessToken(res)
	if err != nil {
		u.log.Error("Error in generating access token", "error", err)
		return nil, err
	}

	expireAt := token.GetExpires()

	u.log.Info("Add user ended", "username", in.Username)

	return &entity.LogInRes{Token: accessToken, UserId: res.ID, ExpireAt: expireAt}, nil
}
func (u *UserUseCase) IsEmailExists(email string) (bool, error) {
	return u.repo.IsEmailExists(email)
}

func (u *UserUseCase) LogIn(in entity.LogInReq) (*entity.LogInRes, error) {
	u.log.Info("Log in started", "username", in.Username)
	res, err := u.repo.GetUserByUsername(in.Username)
	if err != nil {
		u.log.Error("Error in logging in", "error", err)
		return nil, err
	}

	if !help.CheckPasswordHash(in.Password, res.Password) {
		u.log.Error("Error in logging in", "error", "Invalid username or password")
		return nil, fmt.Errorf("Invalid username or password")
	}

	accessToken, err := token.GenerateAccessToken(res)
	if err != nil {
		u.log.Error("Error in generating access token", "error", err)
		return nil, err
	}

	expireAt := token.GetExpires()

	u.log.Info("Log in ended", "username", in.Username)

	return &entity.LogInRes{
		Token:    accessToken,
		UserId:   res.ID,
		ExpireAt: expireAt,
	}, nil
}
