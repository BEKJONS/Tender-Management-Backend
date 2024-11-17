package email

import (
	"tender_management/internal/entity"
	"tender_management/pkg/postgres"
	"time"
	"tender_management/config"
)

var (
	Config = config.NewConfig()
	db, _  = postgres.Connection(Config)
)

func CreateBidMessage(user string, tenderID string, username string) (string, error) {

	mes := "User " + username + " has made a bid on tender " + tenderID + "\n" + "tap to see more info"
	message := entity.Notification{
		UserID:     user,
		Message:    mes,
		Type:       "accept",
		CreatedAt:  time.Now(),
		RelationID: tenderID,
	}

	query := `
		INSERT INTO notifications (user_id, message, type, created_at, relation_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(query, message.UserID, message.Message, message.Type, message.CreatedAt, message.RelationID)
	if err != nil {
		return "", err
	}
	return message.Message, nil
}


func CreateTenderMessage(user string, tenderID string, username string) (string, error) {

	mes := "User " + username + " has created a tender " + tenderID + "\n" + "tap to see more info"
	message := entity.Notification{
		UserID:     user,
		Message:    mes,
		Type:       "accept",
		CreatedAt:  time.Now(),
		RelationID: tenderID,
	}

	query := `
		INSERT INTO notifications (user_id, message, type, created_at, relation_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(query, message.UserID, message.Message, message.Type, message.CreatedAt, message.RelationID)
	if err != nil {
		return "", err
	}
	return message.Message, nil
}