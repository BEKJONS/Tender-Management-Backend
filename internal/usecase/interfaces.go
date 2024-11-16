package usecase

import "tender_management/internal/entity"

type UsersRepo interface {
	CreateUser(user entity.User) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
}
type TenderRepo interface {
	CreateTender(in entity.TenderRepoReq) (entity.Tender, error)
	GetTender(tenderID string) (entity.Tender, error)
	ListTenders(clientID string) ([]entity.Tender, error)
	UpdateTenderStatus(tender *entity.UpdateTender) (entity.Message, error)
	DeleteTender(tenderID string) (entity.Message, error)
}

type BidRepo interface {
	SubmitBid(bid entity.Bid) (entity.Bid, error)
	GetBids(in entity.ListBidReq) ([]entity.Bid, error)
}
