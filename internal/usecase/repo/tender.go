package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"tender_management/internal/entity"
	"tender_management/internal/usecase"
)

type TenderRepo struct {
	db *sqlx.DB
}

// NewTenderRepo creates a new instance of TenderRepo.
func NewTenderRepo(db *sqlx.DB) usecase.TenderRepo {
	return &TenderRepo{db: db}
}

// CreateTender inserts a new tender record into the database.
func (r *TenderRepo) CreateTender(in entity.TenderRepoReq) (entity.Tender, error) {
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
func (r *TenderRepo) UpdateTenderStatus(tender *entity.UpdateTender) (entity.Message, error) {
	// Базовая проверка на ID
	if tender.Id == "" {
		return entity.Message{}, fmt.Errorf("tender ID is required")
	}

	// Динамическое построение SQL
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if tender.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, tender.Title)
		argIndex++
	}

	if tender.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, tender.Description)
		argIndex++
	}

	if !tender.Deadline.IsZero() {
		updates = append(updates, fmt.Sprintf("deadline = $%d", argIndex))
		args = append(args, tender.Deadline)
		argIndex++
	}

	if tender.Budget > 0 {
		updates = append(updates, fmt.Sprintf("budget = $%d", argIndex))
		args = append(args, tender.Budget)
		argIndex++
	}

	if tender.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, tender.Status)
		argIndex++
	}

	if len(updates) == 0 {
		return entity.Message{}, fmt.Errorf("no fields to update")
	}

	// Добавить условие WHERE для ID
	args = append(args, tender.Id)
	query := fmt.Sprintf(`
		UPDATE tenders
		SET %s
		WHERE id = $%d
		RETURNING id
	`, strings.Join(updates, ", "), argIndex)

	var updatedId string
	err := r.db.QueryRowx(query, args...).Scan(&updatedId)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to update tender: %w", err)
	}

	return entity.Message{
		Message: fmt.Sprintf("Tender with ID %s successfully updated", updatedId),
	}, nil
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

func (r *TenderRepo) GetUserTenders(userID string) ([]entity.Tender, error) {
	// Проверка входных данных
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	// SQL-запрос для получения тендеров
	query := `
        SELECT id, title, description, deadline, budget, status
        FROM tenders
        WHERE client_id = $1
    `

	// Создаем слайс для хранения результата
	var tenders []entity.Tender

	// Выполняем запрос
	rows, err := r.db.Queryx(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenders for user: %w", err)
	}
	defer rows.Close()

	// Проходимся по результатам
	for rows.Next() {
		var tender entity.Tender
		if err := rows.StructScan(&tender); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, tender)
	}

	// Проверяем на ошибки чтения строк
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tenders, nil
}

func (r *TenderRepo) CloseTenders(tenderId string) error {

	_, err := r.db.Exec("UPDATE tenders SET status = 'closed' WHERE id = $1", tenderId)
	if err != nil {
		return fmt.Errorf("failed to close tender: %w", err)
	}

	return nil
}

func (r *TenderRepo) AwardedBide(in *entity.Awarded) (*entity.AwardedRes, error) {
	res := &entity.AwardedRes{}

	_, err := r.db.Exec("UPDATE tenders SET status = 'awarded' WHERE id = $1", in.TenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to update tender: %w", err)
	}

	return res, nil
}
