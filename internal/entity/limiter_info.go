package entity

import (
	"time"

	"github.com/google/uuid"
)

type LimiterInfo struct {
	Id              uuid.UUID `json:"id"`
	IP              string    `json:"ip"`
	TimesRequested  int       `json:"timesRequested"`
	LastRequestDate time.Time `json:"lastRequestDate"`
	Blocked         bool      `json:"blocked"`
}
