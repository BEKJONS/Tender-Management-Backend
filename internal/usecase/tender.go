package usecase

import (
	"errors"
	"log/slog"
	"tender_management/internal/entity"
	"time"
)

type TenderService struct {
	repo TenderRepo
	log  *slog.Logger
}

// NewTenderService creates a new instance of TenderService.
func NewTenderService(repo TenderRepo, log *slog.Logger) *TenderService {
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
	if in.Deadline.Before(time.Now()) {
		s.log.Error("error creating tender", "error", errors.New("deadline must be in the future"))
		return entity.Tender{}, errors.New("deadline must be in the future")
	}
	if in.Budget <= 0 {
		s.log.Error("error creating tender", "error", errors.New("budget must be greater than 0"))
		return entity.Tender{}, errors.New("budget must be greater than 0")
	}

	// Set default status if not provided
	if in.Status == "" {
		in.Status = "open"
	}

	// Call repository to create the tender
	tender, err := s.repo.CreateTender(in)
	if err != nil {
		s.log.Error("error creating tender", "error", err)
		return entity.Tender{}, err
	}
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

func (s *TenderService) UpdateTenderStatus(tenderID, status string) (entity.Message, error) {
	s.log.Info("started updating tender status", "id", tenderID, "status", status)
	defer s.log.Info("ended updating tender status", "id", tenderID, "status", status)

	// Validating the status value
	validStatuses := map[string]bool{"open": true, "closed": true, "awarded": true}
	if !validStatuses[status] {
		s.log.Error("error updating tender status", "error", errors.New("invalid status value"))
		return entity.Message{}, errors.New("invalid status value")
	}

	// Call repository to update status
	msg, err := s.repo.UpdateTenderStatus(tenderID, status)
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
