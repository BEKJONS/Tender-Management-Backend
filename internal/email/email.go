package email

import (
	"bytes"
	"context"
	"log"
	"net/smtp"
	"tender_management/config"

	rdb "tender_management/pkg/redis"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func SendEmail(ctx context.Context, cfg *config.Config, rc *redis.Client, email string, message string) (string, error) {

	err := rdb.Storemessage(ctx, rc, email, message)
	if err != nil {
		return "", err
	}

	_, err = rdb.Getmessage(ctx, rc, email)
	if err != nil {
		return "", err
	}

	err = sendmessage(cfg, email, message)
	if err != nil {
		return "", err
	}

	return message, nil
}

func sendmessage(cfg *config.Config, email string, message string) error {
	from := cfg.EMAIL
	password := cfg.APP_KEY

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer
	body.Write([]byte(message)) // Set the message content in the body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	log.Println("Email sended to:", email)
	return nil
}
