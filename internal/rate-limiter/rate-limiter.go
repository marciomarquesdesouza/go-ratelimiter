package ratelimiter

import (
	"time"

	"github.com/google/uuid"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/entity"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/infra/database"
)

func CheckLimitReached(ip string, limit int, blockingTimeSeconds int, requestRepository database.LimiterInfoRepositoryInterface) (bool, error) {
	ipData, err := requestRepository.GetByIP(ip)
	if err != nil && err.Error() != "IP not found" {
		return false, err
	}

	if ipData == nil {
		ipData = &entity.LimiterInfo{
			Id:              uuid.New(),
			IP:              ip,
			TimesRequested:  1,
			LastRequestDate: time.Now().Local(),
			Blocked:         false,
		}

		requestRepository.Save(ipData)
	} else {
		if ipData.Blocked {
			blockingTimeExpired := time.Now().After(ipData.LastRequestDate.Add(time.Second * time.Duration(blockingTimeSeconds)))

			if blockingTimeExpired {
				ipData.Blocked = false
				ipData.TimesRequested = 1
				ipData.LastRequestDate = time.Now()
			}
		} else {
			if ipData.TimesRequested+1 > limit {
				ipData.Blocked = true
				ipData.LastRequestDate = time.Now()
			} else if time.Now().Before(ipData.LastRequestDate.Add(time.Second)) {
				ipData.TimesRequested++
			} else {
				ipData.TimesRequested = 1
				ipData.LastRequestDate = time.Now()
			}
		}

		requestRepository.Update(ipData)
	}

	return ipData.Blocked, nil
}
