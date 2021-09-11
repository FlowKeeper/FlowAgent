package cache

import "github.com/google/uuid"

//Config stores the current instance config retrieved from the sqlite database
var Config struct {
	AgentUUID   uuid.UUID
	ScraperUUID uuid.UUID
}
