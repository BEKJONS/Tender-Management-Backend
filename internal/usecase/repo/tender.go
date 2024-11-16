package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"tender_management/internal/entity"
)

type TenderRepo struct {
	db *sqlx.DB
}

// NewTenderRepo creates a new instance of TenderRepo.
func NewTenderRepo(db *sqlx.DB) *TenderRepo {
	return &TenderRepo{db: db}
}

// CreateTender inserts a new tender record into the database.
func (r *TenderRepo) CreateTender(in entity.TenderReq) (entity.Tender, error) {
	query := `
		INSERT INTO tenders (client_id, title, description, deadline, budget, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, client_id, title, description, deadline, budget, status
	`
	var tender entity.Tender
	err := r.db.QueryRowx(query, in.ClientID, in.Title, in.Description, in.Deadline, in.Budget, in.Status).
		Scan(&tender.ID, &tender.ClientID, &tender.Title, &tender.Description, &tender.Deadline, &tender.Budget, &tender.Status)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("failed to create tender: %w", err)
	}
	return tender, nil
}

// GetTender retrieves a tender by its ID.
func (r *TenderRepo) GetTender(tenderID string) (entity.Tender, error) {
	var tender entity.Tender
	query := `
		SELECT id, client_id, title, description, deadline, budget, status
		FROM tenders
		WHERE id = $1
	`
	err := r.db.Get(&tender, query, tenderID)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("failed to get tender: %w", err)
	}
	return tender, nil
}

// ListTenders retrieves all tenders, optionally filtered by client ID.
func (r *TenderRepo) ListTenders(clientID string) ([]entity.Tender, error) {
	var tenders []entity.Tender
	query := `
		SELECT id, client_id, title, description, deadline, budget, status
		FROM tenders
	`
	var err error
	if clientID != "" {
		query += " WHERE client_id = $1 ORDER BY deadline DESC"
		err = r.db.Select(&tenders, query, clientID)
	} else {
		query += " ORDER BY deadline DESC"
		err = r.db.Select(&tenders, query)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list tenders: %w", err)
	}
	return tenders, nil
}

// UpdateTenderStatus updates the status of a tender.
func (r *TenderRepo) UpdateTenderStatus(tenderID string, status string) (entity.Message, error) {
	query := `
		UPDATE tenders
		SET status = $1
		WHERE id = $2
		RETURNING id
	`
	var id string
	err := r.db.QueryRow(query, status, tenderID).Scan(&id)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to update tender status: %w", err)
	}
	return entity.Message{Message: "Tender status updated successfully"}, nil
}

// DeleteTender removes a tender by its ID.
func (r *TenderRepo) DeleteTender(tenderID string) (entity.Message, error) {
	query := `
		DELETE FROM tenders
		WHERE id = $1
	`
	res, err := r.db.Exec(query, tenderID)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to delete tender: %w", err)
	}
	rows, _ := res.RowsAffected()
	return entity.Message{Message: fmt.Sprintf("Deleted %d tender(s)", rows)}, nil
}
