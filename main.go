package main

import (
	"fmt"

	"github.com/FlowKeeper/FlowAgent/v2/cache"
	"github.com/FlowKeeper/FlowAgent/v2/config"
	"github.com/FlowKeeper/FlowAgent/v2/fetcher"
	"github.com/FlowKeeper/FlowAgent/v2/webserver"
	"github.com/FlowKeeper/FlowUtils/v2/flowutils"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
)

func main() {
	logger.Info("MAIN", "Starting FlowKeeper FlowAgent")
	utilsVersion := flowutils.Version()
	logger.Info("UTILS", fmt.Sprintf("Running FlowUtils Version: %d-%d-%s", utilsVersion.Major, utilsVersion.Minor, utilsVersion.Comment))

	config.Init()
	cache.Init()
	go fetcher.Init()
	webserver.Init()
}
