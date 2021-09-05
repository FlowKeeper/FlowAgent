package cache

import "github.com/google/uuid"

type SampleConfig struct {
	AgentUUID   uuid.UUID
	ScraperUUID uuid.UUID
}

var Config SampleConfig
