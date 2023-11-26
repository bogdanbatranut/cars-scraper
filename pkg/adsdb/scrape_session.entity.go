package adsdb

import (
	"github.com/google/uuid"
)

type ScrapeSession struct {
	MarketID   uint
	CriteriaID uint
	UUID       uuid.UUID
}
