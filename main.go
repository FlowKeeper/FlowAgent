package main

import (
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/cache"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/config"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/fetcher"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/webserver"
)

func main() {
	logger.Info("MAIN", "Starting FlowKeeper FlowAgent")
	config.Init()
	cache.Init()
	go fetcher.Init()
	webserver.Init()
}
