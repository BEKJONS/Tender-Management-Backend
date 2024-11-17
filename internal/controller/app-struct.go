package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"tender_management/internal/usecase"
	"tender_management/internal/usecase/cashing"
	"tender_management/internal/usecase/repo"
)

type Controller struct {
	Auth *usecase.UserUseCase
	Tend *usecase.TenderService
	Bid  *usecase.BidService
}

func NewController(db *sqlx.DB, log *slog.Logger, rd *redis.Client) *Controller {

	authRepo := repo.NewUserRepo(db)
	tendRepo := repo.NewTenderRepo(db)
	bidRepo := repo.NewBidRepo(db)
	cash := cashing.NewTenderCash(log, rd)

	ctr := &Controller{
		Auth: usecase.NewUserUseCase(authRepo, log),
		Tend: usecase.NewTenderService(tendRepo, bidRepo, cash, log),
		Bid:  usecase.NewBidUseCase(bidRepo, log),
	}

	return ctr
}
