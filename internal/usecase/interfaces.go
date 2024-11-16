package usecase

import "tender_management/internal/entity"

type UsersRepo interface {
	CreateUser(user entity.User) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
}
type TenderRepo interface {
	CreateTender(in entity.TenderReq) (entity.Tender, error)
	GetTender(tenderID string) (entity.Tender, error)
	ListTenders(clientID string) ([]entity.Tender, error)
	UpdateTenderStatus(tenderID string, status string) (entity.Message, error)
	DeleteTender(tenderID string) (entity.Message, error)
}
