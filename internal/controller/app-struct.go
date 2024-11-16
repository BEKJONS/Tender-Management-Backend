package controller

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"tender_management/internal/usecase"
	"tender_management/internal/usecase/repo"
)

type Controller struct {
	Auth *usecase.UserUseCase
	Tend *usecase.TenderService
}

func NewController(db *sqlx.DB, log *slog.Logger) *Controller {

	authRepo := repo.NewUserRepo(db)
	tendRepo := repo.NewTenderRepo(db)

	ctr := &Controller{
		Auth: usecase.NewUserUseCase(authRepo, log),
		Tend: usecase.NewTenderService(tendRepo, log),
	}

	return ctr
}
