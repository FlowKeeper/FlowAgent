package config

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

type SampleConfig struct {
	ListenAddress string
	ListenPort    int
	CachePath     string
}

var Config SampleConfig

const linuxConfigPath = "/etc/flowagent/config.json"
const windowsConfigPath = `C:\flowagent\config.json`
const windowsCacheExamplePath = `C:\flowagent\cache.sqlite`
const linuxCacheExamplePath = "/var/lib/flowagent/cache.sqlite"

var configPath string
var cacheExamplePath string

const loggingArea = "CONFIG"

//Init initializes the config struct
func Init() {
	switch runtime.GOOS {
	case "windows":
		configPath = windowsConfigPath
		cacheExamplePath = windowsCacheExamplePath
	case "linux":
		configPath = linuxConfigPath
		cacheExamplePath = linuxCacheExamplePath
	default:
		logger.Fatal(loggingArea, "Agent is running on unsupported platform:", runtime.GOOS)
	}

	content, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		logger.Fatal(loggingArea, fmt.Sprintf("No such file or directory: %s. Was the agent intalled properly?", configPath))
	}

	if err := json.Unmarshal(content, &Config); err != nil {
		logger.Fatal(loggingArea, "Couldn't parse config:", err)
	}

	//Sanity check config
	if stringHelper.IsEmpty(Config.ListenAddress) {
		logger.Fatal(loggingArea, "ListenAddress is malformed. Example: \"0.0.0.0\"")
	}

	if Config.ListenPort == 0 {
		logger.Fatal(loggingArea, "ListenPort is malformed. Example: 5000")
	}

	if stringHelper.IsEmpty(Config.CachePath) {
		logger.Fatal(loggingArea, "CachePath is malformed. Example:", cacheExamplePath)
	}

	logger.Info(loggingArea, "Config is operational")
}
