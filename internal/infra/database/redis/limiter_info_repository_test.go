package redis

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestSaveLimiterInfo(t *testing.T) {
	repository := NewLimiterInfoRepository("localhost:6379", "", 0)

	newLimiterInfo := &entity.LimiterInfo{
		Id:              uuid.New(),
		IP:              "192.168.0.1",
		TimesRequested:  1,
		LastRequestDate: time.Now().Local(),
		Blocked:         false,
	}

	err := repository.Save(newLimiterInfo)
	assert.Nil(t, err)

	ipData, err := repository.GetByIP("192.168.0.1")
	assert.Nil(t, err)

	assert.Equal(t, ipData.IP, newLimiterInfo.IP)
	assert.Equal(t, ipData.TimesRequested, newLimiterInfo.TimesRequested)
	assert.Equal(t, ipData.LastRequestDate, newLimiterInfo.LastRequestDate)
	assert.Equal(t, ipData.Blocked, newLimiterInfo.Blocked)
}

func TestUpdateLimiterInfoSuccess(t *testing.T) {
	repository := NewLimiterInfoRepository("localhost:6379", "", 0)

	newLimiterInfo := &entity.LimiterInfo{
		Id:              uuid.New(),
		IP:              "192.168.0.1",
		TimesRequested:  1,
		LastRequestDate: time.Now().Local(),
		Blocked:         false,
	}

	err := repository.Save(newLimiterInfo)
	assert.Nil(t, err)

	err = repository.Update(newLimiterInfo)
	assert.Nil(t, err)

	ipData, err := repository.GetByIP("192.168.0.1")
	assert.Nil(t, err)

	assert.Equal(t, ipData.IP, newLimiterInfo.IP)
	assert.Equal(t, ipData.TimesRequested, newLimiterInfo.TimesRequested)
	assert.Equal(t, ipData.LastRequestDate, newLimiterInfo.LastRequestDate)
	assert.Equal(t, ipData.Blocked, newLimiterInfo.Blocked)
}

func TestUpdateLimiterInfoNotExists(t *testing.T) {
	repository := NewLimiterInfoRepository("localhost:6379", "", 0)

	limiterInfo := &entity.LimiterInfo{
		Id:              uuid.New(),
		IP:              "192.168.3.1",
		TimesRequested:  1,
		LastRequestDate: time.Now().Local(),
		Blocked:         false,
	}

	err := repository.Update(limiterInfo)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "IP not found")
}
