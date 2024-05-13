package database

import "github.com/marciomarquesdesouza/go-rate-limiter/internal/entity"

type LimiterInfoRepositoryInterface interface {
	Save(r *entity.LimiterInfo) error
	GetByIP(ip string) (*entity.LimiterInfo, error)
	Update(r *entity.LimiterInfo) error
}
