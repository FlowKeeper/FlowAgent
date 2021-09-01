package cache

import "github.com/google/uuid"

type SampleConfig struct {
	AgentID   uuid.UUID
	ScraperID uuid.UUID
}

var Config SampleConfig
