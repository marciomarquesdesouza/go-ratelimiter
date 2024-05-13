package ratelimiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/marciomarquesdesouza/go-rate-limiter/internal/infra/database/redis"
	"github.com/stretchr/testify/assert"
)

func TestIPUnblocked(t *testing.T) {
	repository := redis.NewLimiterInfoRepository("localhost:6379", "", 0)

	//First request
	limitReached, err := CheckLimitReached("192.168.1.2", 4, 10, repository)
	assert.Nil(t, err)
	assert.False(t, limitReached)

	ipData, err := repository.GetByIP("192.168.1.2")
	assert.Nil(t, err)

	// assert.False(t, ipData.Blocked)
	assert.Equal(t, ipData.IP, "192.168.1.2")
	assert.Equal(t, ipData.Blocked, false)
}

func TestRateLimiter(t *testing.T) {
	repository := redis.NewLimiterInfoRepository("localhost:6379", "", 0)

	//First request
	limitReached, err := CheckLimitReached("192.168.1.2", 3, 60, repository)
	assert.Nil(t, err)
	assert.False(t, limitReached)

	ipData, err := repository.GetByIP("192.168.1.2")
	assert.Nil(t, err)

	assert.Equal(t, ipData.IP, "192.168.1.2")
	assert.Equal(t, ipData.Blocked, false)

	//Second request
	limitReached, err = CheckLimitReached(ipData.IP, 3, 60, repository)
	assert.Nil(t, err)
	assert.False(t, limitReached)

	//Third request
	limitReached, err = CheckLimitReached(ipData.IP, 4, 60, repository)
	assert.Nil(t, err)
	assert.False(t, limitReached)

	//Fourth request -> blocked
	limitReached, err = CheckLimitReached(ipData.IP, 3, 60, repository)
	assert.Nil(t, err)
	assert.True(t, limitReached)

	ipData, err = repository.GetByIP("192.168.1.2")
	assert.Nil(t, err)

	assert.True(t, ipData.Blocked)
}

func TestIPBlocked(t *testing.T) {
	repository := redis.NewLimiterInfoRepository("localhost:6379", "", 0)

	//First request
	limitReached, err := CheckLimitReached("192.168.1.2", 1, 10, repository)
	assert.Nil(t, err)
	assert.True(t, limitReached)

	ipData, err := repository.GetByIP("192.168.1.2")
	assert.Nil(t, err)

	assert.True(t, ipData.Blocked)
}

func BenchmarkLimitReached(b *testing.B) {
	repository := redis.NewLimiterInfoRepository("localhost:6379", "", 0)

	for i := 0; i < b.N; i++ {
		limitReached, err := CheckLimitReached("192.168.1.2", 3, 10, repository)
		if err != nil {
			fmt.Println(err.Error())
		}

		if limitReached {
			fmt.Println("Limit reached. Waiting for the blocking time.")
			time.Sleep(time.Second)
		} else {
			fmt.Println("You are able to call the API.")
		}
	}
}
