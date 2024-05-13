package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/marciomarquesdesouza/go-rate-limiter/internal/entity"
	"github.com/redis/go-redis/v9"
)

type LimiterInfoRepository struct {
	Client *redis.Client
}

func NewLimiterInfoRepository(address string, password string, db int) *LimiterInfoRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &LimiterInfoRepository{
		Client: client,
	}
}

func (r *LimiterInfoRepository) Save(limiterInfo *entity.LimiterInfo) error {
	json, err := json.Marshal(limiterInfo)
	if err != nil {
		return err
	}

	err = r.Client.Set(context.Background(), limiterInfo.IP, json, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *LimiterInfoRepository) GetByIP(ip string) (*entity.LimiterInfo, error) {
	data, err := r.Client.Get(context.Background(), ip).Result()
	if err == redis.Nil {
		return nil, errors.New("IP not found")
	} else if err != nil {
		return nil, err
	}

	var limiterInfo entity.LimiterInfo

	json.Unmarshal([]byte(data), &limiterInfo)

	return &limiterInfo, nil
}

func (r *LimiterInfoRepository) Update(limiterInfo *entity.LimiterInfo) error {
	_, err := r.GetByIP(limiterInfo.IP)
	if err != nil {
		return err
	}

	err = r.Save(limiterInfo)
	if err != nil {
		return err
	}

	return nil
}
