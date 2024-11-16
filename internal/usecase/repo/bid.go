package repo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"tender_management/internal/entity"
)

type BidRepo struct {
	db *sqlx.DB
}

func NewBidRepo(db *sqlx.DB) *BidRepo {
	return &BidRepo{db: db}
}

func (r *BidRepo) SubmitBid(bid entity.Bid) (entity.Bid, error) {
	// Generate a new UUID for the bid
	bidID := uuid.New().String()

	// Insert the bid into the database
	query := `INSERT INTO bids (id, tender_id, contractor_id, price, delivery_time, comments, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.Get(&bid.ID, query, bidID, bid.TenderID, bid.ContractorID, bid.Price, bid.DeliveryTime, bid.Comments, bid.Status)
	if err != nil {
		return entity.Bid{}, err
	}

	// Return the newly created bid
	return bid, nil
}

// GetBids retrieves bids for a specific tender, with optional filters for price and delivery time, and sorting.
func (r *BidRepo) GetBids(in entity.ListBidReq) ([]entity.Bid, error) {
	var bids []entity.Bid
	var query = `SELECT id, tender_id, contractor_id, price, delivery_time, comments, status 
                 FROM bids WHERE tender_id = $1`

	// Add price filter if provided
	if in.PriceFilter != nil {
		query += fmt.Sprintf(" AND price <= $2")
	}

	// Add delivery_time filter if provided
	if in.DeliveryTimeFilter != nil {
		query += fmt.Sprintf(" AND delivery_time <= $3")
	}

	// Add comments filter if provided
	if in.Comments != "" {
		query += fmt.Sprintf(" AND comments ILIKE '%%%s%%'", in.Comments)
	}

	// Add status filter if provided
	if in.Status != "" {
		query += fmt.Sprintf(" AND status = '%s'", in.Status)
	}
	if in.ClientID != "" {
		query += fmt.Sprintf(" AND tender_id IN (SELECT id FROM tenders WHERE client_id = '%s')", in.ClientID)
	}

	// Prepare arguments based on the filters
	args := []interface{}{in.TenderID}
	if in.PriceFilter != nil {
		args = append(args, *in.PriceFilter)
	}
	if in.DeliveryTimeFilter != nil {
		args = append(args, *in.DeliveryTimeFilter)
	}

	err := r.db.Select(&bids, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bids: %w", err)
	}

	return bids, nil
}

func (r *BidRepo) GetUserBids(userID string) ([]entity.Bid, error) {
	// Проверка входных данных
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	// SQL-запрос для получения ставок пользователя
	query := `
        SELECT id, tender_id, contractor_id, price, delivery_time, comments, status
        FROM bids
        WHERE contractor_id = $1
    `

	// Создаем слайс для хранения результата
	var bids []entity.Bid

	// Выполняем запрос
	rows, err := r.db.Queryx(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bids for user: %w", err)
	}
	defer rows.Close()

	// Проходимся по результатам
	for rows.Next() {
		var bid entity.Bid
		if err := rows.StructScan(&bid); err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	// Проверяем на ошибки чтения строк
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return bids, nil
}
