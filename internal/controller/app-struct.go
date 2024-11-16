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
	Bid  *usecase.BidService
}

func NewController(db *sqlx.DB, log *slog.Logger) *Controller {

	authRepo := repo.NewUserRepo(db)
	tendRepo := repo.NewTenderRepo(db)
	bidRepo := repo.NewBidRepo(db)

	ctr := &Controller{
		Auth: usecase.NewUserUseCase(authRepo, log),
		Tend: usecase.NewTenderService(tendRepo, log),
		Bid:  usecase.NewBidUseCase(bidRepo, log),
	}

	return ctr
}
