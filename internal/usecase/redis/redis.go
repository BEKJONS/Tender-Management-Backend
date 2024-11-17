package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis-сервера
		Password: "",               // Пароль (если отсутствует, оставьте пустым)
		DB:       0,                // Номер базы данных (по умолчанию 0)
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	return rdb
}
