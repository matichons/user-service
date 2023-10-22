package repository

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type userCacheRepository struct {
	redisClient *redis.Client
}

const (
	basePrefix = "user"
)

func NewCacheRepository(redisClient *redis.Client) *userCacheRepository {
	return &userCacheRepository{redisClient: redisClient}
}

func (u *userCacheRepository) Incr(ctx *gin.Context, key string) (int64, error) {
	keyCache := u.getKeyWithPrefix("username", key)
	val, err := u.redisClient.Incr(ctx, keyCache).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing key: %v", err)
	}

	if err := u.setExpire(ctx, keyCache); err != nil {
		return val, err
	}
	return val, nil
}

func (u *userCacheRepository) GetByUsernameLogin(ctx *gin.Context, key string) (int64, error) {
	keyCache := u.getKeyWithPrefix("username", key)
	val, err := u.redisClient.Get(ctx, keyCache).Int64()
	if err != nil {
		return 0, fmt.Errorf("error incrementing key: %v", err)
	}

	return val, nil
}

func (u *userCacheRepository) Delete(ctx *gin.Context, key string) error {
	keyCache := u.getKeyWithPrefix("username", key)

	if err := u.redisClient.Del(ctx, keyCache); err != nil {
		return err.Err()
	}

	return nil
}
func (u *userCacheRepository) setExpire(ctx *gin.Context, key string) error {
	if err := u.redisClient.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("error setting TTL: %v", err)
	}
	return nil
}

func (u *userCacheRepository) getKeyWithPrefix(category string, key string) string {
	return fmt.Sprintf("%s-%s: %s", basePrefix, category, key)
}
