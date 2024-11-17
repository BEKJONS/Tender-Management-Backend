package cashing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"tender_management/internal/entity"
)

type TenderCash struct {
	log   *slog.Logger
	redis *redis.Client
}

func NewTenderCash(log *slog.Logger, redis *redis.Client) *TenderCash {
	return &TenderCash{log: log, redis: redis}
}

func (t *TenderCash) SaveNewTender(in *entity.Tender) error {
	if in == nil {
		t.log.Error("Tender is nil")
		return fmt.Errorf("tender is nil")
	}

	cacheKey := "tenders"

	data, err := json.Marshal(in)
	if err != nil {
		t.log.Error("Failed to marshal tender", "error", err)
		return fmt.Errorf("failed to serialize tender: %w", err)
	}

	err = t.redis.RPush(context.Background(), cacheKey, data).Err()
	if err != nil {
		t.log.Error("Failed to append tender to Redis list", "error", err)
		return fmt.Errorf("failed to append tender to Redis list: %w", err)
	}

	t.log.Info("Tender successfully added to Redis list", "tenderID", in.ID)

	err = t.SaveWithClientID(in)
	if err != nil {
		t.log.Error("Failed to save tender", "error", err)
		return fmt.Errorf("failed to save tender: %w", err)
	}

	return nil
}

func (t *TenderCash) SaveWithClientID(in *entity.Tender) error {
	if in == nil {
		t.log.Error("Tender is nil")
		return fmt.Errorf("tender is nil")
	}

	cacheKey := in.ClientID

	data, err := json.Marshal(in)
	if err != nil {
		t.log.Error("Failed to marshal tender", "error", err)
		return fmt.Errorf("failed to serialize tender: %w", err)
	}

	err = t.redis.RPush(context.Background(), cacheKey, data).Err()
	if err != nil {
		t.log.Error("Failed to append tender to Redis list", "error", err)
		return fmt.Errorf("failed to append tender to Redis list: %w", err)
	}

	t.log.Info("Tender successfully added to Redis list", "tenderID", in.ID)
	return nil
}

func (t *TenderCash) UpdateTender(update *entity.UpdateTender) error {

	var clientID string

	if update.Id == "" {
		t.log.Error("Tender ID is required")
		return fmt.Errorf("tender ID is required")
	}

	cacheKey := "tenders"

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	updated := false
	for i, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}

		if tender.ID == update.Id {
			if update.Title != "" {
				tender.Title = update.Title
			}
			if update.Description != "" {
				tender.Description = update.Description
			}
			if !update.Deadline.IsZero() {
				tender.Deadline = update.Deadline
			}
			if update.Budget > 0 {
				tender.Budget = update.Budget
			}
			if update.Status != "" {
				tender.Status = update.Status
			}

			newData, err := json.Marshal(tender)
			if err != nil {
				t.log.Error("Failed to marshal updated tender", "error", err)
				return fmt.Errorf("failed to serialize updated tender: %w", err)
			}

			err = t.redis.LSet(context.Background(), cacheKey, int64(i), newData).Err()
			if err != nil {
				t.log.Error("Failed to update tender in Redis", "error", err)
				return fmt.Errorf("failed to update tender in Redis: %w", err)
			}

			updated = true
			clientID = tender.ID
			break
		}
	}

	if !updated {
		t.log.Error("Tender not found for update", "tenderID", update.Id)
		return fmt.Errorf("tender with ID %s not found", update.Id)
	}

	t.log.Info("Tender updated successfully", "tenderID", update.Id)

	go func() {
		err := t.UpdateTenderClient(update, clientID)
		if err != nil {
			t.log.Error("Failed to update client in Redis", "error", err)
			return
		}
	}()

	return nil
}

func (t *TenderCash) UpdateTenderClient(update *entity.UpdateTender, clientID string) error {
	if update.Id == "" {
		t.log.Error("Tender ID is required")
		return fmt.Errorf("tender ID is required")
	}

	cacheKey := clientID

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	updated := false
	for i, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}

		if tender.ID == update.Id {
			if update.Title != "" {
				tender.Title = update.Title
			}
			if update.Description != "" {
				tender.Description = update.Description
			}
			if !update.Deadline.IsZero() {
				tender.Deadline = update.Deadline
			}
			if update.Budget > 0 {
				tender.Budget = update.Budget
			}
			if update.Status != "" {
				tender.Status = update.Status
			}

			newData, err := json.Marshal(tender)
			if err != nil {
				t.log.Error("Failed to marshal updated tender", "error", err)
				return fmt.Errorf("failed to serialize updated tender: %w", err)
			}

			err = t.redis.LSet(context.Background(), cacheKey, int64(i), newData).Err()
			if err != nil {
				t.log.Error("Failed to update tender in Redis", "error", err)
				return fmt.Errorf("failed to update tender in Redis: %w", err)
			}

			updated = true
			break
		}
	}

	if !updated {
		t.log.Error("Tender not found for update", "tenderID", update.Id)
		return fmt.Errorf("tender with ID %s not found", update.Id)
	}

	t.log.Info("Tender updated successfully", "tenderID", update.Id)
	return nil
}

func (t *TenderCash) DeleteTender(tenderID string) error {

	var clientID string

	if tenderID == "" {
		t.log.Error("Tender ID is required")
		return fmt.Errorf("tender ID is required")
	}

	cacheKey := "tenders"

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	var updatedList []string
	found := false
	for _, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}

		if tender.ID != tenderID {
			updatedList = append(updatedList, item)
		} else {
			found = true
			clientID = tender.ID
		}
	}

	if !found {
		t.log.Error("Tender not found for deletion", "tenderID", tenderID)
		return fmt.Errorf("tender with ID %s not found", tenderID)
	}

	if err != nil {
		t.log.Error("Failed to delete old tenders list from Redis", "error", err)
		return fmt.Errorf("failed to delete old tenders list from Redis: %w", err)
	}

	if len(updatedList) > 0 {
		if err != nil {
			t.log.Error("Failed to update tenders list in Redis", "error", err)
			return fmt.Errorf("failed to update tenders list in Redis: %w", err)
		}
	}

	t.log.Info("Tender deleted successfully", "tenderID", tenderID)

	go func() {
		err := t.DeleteTenderClient(tenderID, clientID)
		if err != nil {
			t.log.Error("Failed to delete client in Redis", "error", err)
			return
		}
	}()

	return nil
}

func (t *TenderCash) DeleteTenderClient(tenderID, clientID string) error {
	if tenderID == "" {
		t.log.Error("Tender ID is required")
		return fmt.Errorf("tender ID is required")
	}

	cacheKey := clientID

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	var updatedList []string
	found := false
	for _, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}

		if tender.ID != tenderID {
			updatedList = append(updatedList, item)
		} else {
			found = true
		}
	}

	if !found {
		t.log.Error("Tender not found for deletion", "tenderID", tenderID)
		return fmt.Errorf("tender with ID %s not found", tenderID)
	}

	if err != nil {
		t.log.Error("Failed to delete old tenders list from Redis", "error", err)
		return fmt.Errorf("failed to delete old tenders list from Redis: %w", err)
	}

	if len(updatedList) > 0 {
		if err != nil {
			t.log.Error("Failed to update tenders list in Redis", "error", err)
			return fmt.Errorf("failed to update tenders list in Redis: %w", err)
		}
	}

	t.log.Info("Tender deleted successfully", "tenderID", tenderID)
	return nil
}

func (t *TenderCash) GetAllTenders() ([]entity.Tender, error) {
	cacheKey := "tenders"

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return nil, fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	var tenders []entity.Tender
	for _, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}
		tenders = append(tenders, tender)
	}

	t.log.Info("Successfully retrieved tenders from Redis", "count", len(tenders))
	return tenders, nil
}

func (t *TenderCash) GetUserTenders(clientID string) ([]entity.Tender, error) {
	cacheKey := clientID

	data, err := t.redis.LRange(context.Background(), cacheKey, 0, -1).Result()
	if err != nil {
		t.log.Error("Failed to retrieve tenders from Redis", "error", err)
		return nil, fmt.Errorf("failed to retrieve tenders from Redis: %w", err)
	}

	var tenders []entity.Tender
	for _, item := range data {
		var tender entity.Tender
		if err := json.Unmarshal([]byte(item), &tender); err != nil {
			t.log.Error("Failed to unmarshal tender", "error", err)
			continue
		}
		tenders = append(tenders, tender)
	}

	t.log.Info("Successfully retrieved tenders from Redis", "count", len(tenders))
	return tenders, nil
}
