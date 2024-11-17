package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func Storemessage(ctx context.Context, rdb *redis.Client, email, message string) error {
	err := rdb.Set(ctx, "miniTwitter:"+email, message, time.Minute*3).Err()
	if err != nil {
		return errors.Wrap(err, "failed to store message in redis")
	}

	return nil
}

func Getmessage(ctx context.Context, rdb *redis.Client, email string) (string, error) {
	message, err := rdb.Get(ctx, "miniTwitter:"+email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("message not found for " + email)
		}
		return "", errors.Wrap(err, "failed to get message from redis")
	}

	return message, nil
}

func Deletemessage(ctx context.Context, rdb *redis.Client, email string) error {
	err := rdb.Del(ctx, "miniTwitter:"+email).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete message from redis")
	}

	return nil
}
