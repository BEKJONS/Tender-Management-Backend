package rate_limiting

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redis  *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(redis *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{redis: redis, limit: limit, window: window}
}

// Проверка лимита для клиента
func (r *RateLimiter) Allow(clientID string) (bool, error) {
	key := fmt.Sprintf("rate:%s", clientID)

	// Увеличиваем счётчик запросов
	count, err := r.redis.Incr(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	// Если это первый запрос, устанавливаем срок действия ключа
	if count == 1 {
		r.redis.Expire(context.Background(), key, r.window)
	}

	// Проверяем, не превышен ли лимит
	if int(count) > r.limit {
		return false, nil
	}

	return true, nil
}
