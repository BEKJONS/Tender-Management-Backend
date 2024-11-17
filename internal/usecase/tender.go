package usecase

import (
	"context"
	"errors"
	"log/slog"
	"tender_management/internal/entity"
	"time"
)

type TenderService struct {
	repo TenderRepo
	bid  BidRepo
	log  *slog.Logger
}

// NewTenderService creates a new instance of TenderService.
func NewTenderService(repo TenderRepo, bid BidRepo, log *slog.Logger) *TenderService {
	return &TenderService{repo: repo, log: log}
}

func (s *TenderService) CreateTender(in entity.TenderReq) (entity.Tender, error) {
	s.log.Info("started creating tender", "title", in.Title)

	defer s.log.Info("ended creating tender", "title", in.Title)

	// Validation checks
	if in.Title == "" {
		s.log.Error("error creating tender", "error", errors.New("title is required"))
		return entity.Tender{}, errors.New("title is required")
	}

	if in.Deadline.After(time.Now()) {
		s.log.Error("error creating tender", "error", errors.New("deadline must be in the future"))
		return entity.Tender{}, errors.New("deadline must be in the future")
	}

	if in.Budget <= 0 {
		s.log.Error("error creating tender", "error", errors.New("budget must be greater than 0"))
		return entity.Tender{}, errors.New("budget must be greater than 0")
	}

	req := entity.TenderRepoReq{
		ClientID:    in.ClientID,
		Title:       in.Title,
		Description: in.Description,
		Deadline:    in.Deadline,
		Budget:      in.Budget,
		Status:      "open",
	}

	tender, err := s.repo.CreateTender(req)
	if err != nil {
		s.log.Error("error creating tender", "error", err)
		return entity.Tender{}, err
	}

	ctx, _ := context.WithDeadline(context.Background(), req.Deadline)

	go func() {
		select {
		case <-ctx.Done():
			err := s.repo.CloseTenders(tender.ID)
			if err != nil {
				s.log.Error("error closing tender", "error", err)
				return
			}
		}
	}()

	return tender, nil
}

func (s *TenderService) GetTender(tenderID string) (entity.Tender, error) {
	s.log.Info("started getting tender", "id", tenderID)
	defer s.log.Info("ended getting tender", "id", tenderID)

	tender, err := s.repo.GetTender(tenderID)
	if err != nil {
		s.log.Error("error getting tender", "error", err)
		return entity.Tender{}, err
	}
	return tender, nil
}

func (s *TenderService) ListTenders(clientID string) ([]entity.Tender, error) {
	s.log.Info("started listing tenders", "clientID", clientID)
	defer s.log.Info("ended listing tenders", "clientID", clientID)

	tenders, err := s.repo.ListTenders(clientID)
	if err != nil {
		s.log.Error("error listing tenders", "error", err)
		return nil, err
	}
	return tenders, nil
}

func (s *TenderService) UpdateTenderStatus(req *entity.UpdateTender) (entity.Message, error) {
	s.log.Info("started updating tender status", "id", req.Id, "status", req.Status)
	defer s.log.Info("ended updating tender status", "id", req.Id, "status", req.Status)

	// Validating the status value
	validStatuses := map[string]bool{"open": true, "closed": true, "awarded": true}
	if !validStatuses[req.Status] {
		s.log.Error("error updating tender status", "error", errors.New("invalid status value"))
		return entity.Message{}, errors.New("invalid status value")
	}

	// Call repository to update status
	msg, err := s.repo.UpdateTenderStatus(req)
	if err != nil {
		s.log.Error("error updating tender status", "error", err)
		return entity.Message{}, err
	}
	return msg, nil
}

func (s *TenderService) DeleteTender(tenderID string) (entity.Message, error) {
	s.log.Info("started deleting tender", "id", tenderID)
	defer s.log.Info("ended deleting tender", "id", tenderID)

	msg, err := s.repo.DeleteTender(tenderID)
	if err != nil {
		s.log.Error("error deleting tender", "error", err)
		return entity.Message{}, err
	}
	return msg, nil
}

func (s *TenderService) GetUserTenders(clientID string) ([]entity.Tender, error) {
	s.log.Info("started getting user tenders", "clientID", clientID)
	defer s.log.Info("ended getting user tenders", "clientID", clientID)

	res, err := s.repo.GetUserTenders(clientID)
	if err != nil {
		s.log.Error("error getting user tenders", "error", err)
		return nil, err
	}

	return res, nil
}

func (s *TenderService) AwardTender(in *entity.Awarded) (*entity.AwardedRes, error) {

	res, err := s.repo.AwardedBide(in)
	if err != nil {
		s.log.Error("error awarded tender", "error", err)
		return nil, err
	}

	res, err = s.bid.AwardedBide(in)
	if err != nil {
		s.log.Error("error awarded tender", "error", err)
		return nil, err
	}

	return res, nil
}
