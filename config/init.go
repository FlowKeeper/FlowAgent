package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

//Config stores the current configuration of the agent
var Config struct {
	Name               string
	Debug              bool
	ListenAddress      string
	ListenPort         int
	CachePath          string
	ServerAddress      string
	ServerAdressParsed url.URL
}

const linuxConfigPath = "/etc/flowagent/config.json"
const windowsConfigPath = `C:\flowagent\config.json`
const windowsCacheExamplePath = `C:\flowagent\cache.sqlite`
const linuxCacheExamplePath = "/var/lib/flowagent/cache.sqlite"

var configPath string
var cacheExamplePath string

const loggingArea = "CONFIG"

//Init initializes the config struct
func Init() {
	var useENV bool
	var err error
	if !stringHelper.IsEmpty(os.Getenv("FLOW_ENV")) {
		useENV, err = strconv.ParseBool(os.Getenv("FLOW_ENV"))
	}

	//Determine OS and use specific paths for config and cache
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

	if !useENV || err != nil {
		content, err := os.ReadFile(configPath)
		if os.IsNotExist(err) {
			logger.Fatal(loggingArea, fmt.Sprintf("No such file or directory: %s. Was the agent intalled properly?", configPath))
		}

		if err := json.Unmarshal(content, &Config); err != nil {
			logger.Fatal(loggingArea, "Couldn't parse config:", err)
		}
	} else {
		readEnv()
	}

	//Check if all fields are set correctly
	if stringHelper.IsEmpty(Config.ListenAddress) {
		if !useENV {
			logger.Fatal(loggingArea, "ListenAddress is malformed. Example: \"0.0.0.0\"")
		} else {
			logger.Warning(loggingArea, "Using default ListenAddress: 0.0.0.0")
			Config.ListenAddress = "0.0.0.0"
		}
	}

	if Config.ListenPort == 0 {
		if !useENV {
			logger.Fatal(loggingArea, "ListenPort is malformed. Example: 5000")
		} else {
			logger.Warning(loggingArea, "Using default ListenPort: 5000")
			Config.ListenPort = 5000
		}
	}

	if stringHelper.IsEmpty(Config.CachePath) {
		if !useENV {
			logger.Fatal(loggingArea, "CachePath is malformed. Example:", cacheExamplePath)
		} else {
			logger.Warning(loggingArea, "Using default CachePath:", cacheExamplePath)
			Config.CachePath = cacheExamplePath
		}
	}

	if stringHelper.IsEmpty(Config.ServerAddress) {
		if !useENV {
			logger.Fatal(loggingArea, "ServerAdress is malformed. Example: http://my-server.domain.tld:5000")
		} else {
			logger.Fatal(loggingArea, "ServerAddress is malformed / not set. Please set the ENV Variable: FLOW_SERVER_ADDRESS")
		}
	}

	if stringHelper.IsEmpty(Config.Name) {
		if !useENV {
			logger.Fatal(loggingArea, "Name is malformed. Example: MyAgent")
		} else {
			hostname, _ := os.Hostname()
			logger.Warning(loggingArea, "Using hostname as agent's name:", hostname)
			Config.Name = hostname
		}
	}

	//Parse ServerAddress for easier use later
	url, err := url.Parse(Config.ServerAddress)
	if err != nil {
		logger.Fatal(loggingArea, "Couldn't parse ServerAddress:", err)
	}

	Config.ServerAdressParsed = *url

	if Config.Debug {
		logger.EnableDebugLog()
		logger.Debug(loggingArea, "Enabled debug logging")
	}

	logger.Info(loggingArea, "Config is operational")
}

func readEnv() {
	logger.Info(loggingArea, "Using ENV Variables as config")

	Config.ListenAddress = os.Getenv("FLOW_LISTEN_ADDRESS")
	Config.CachePath = os.Getenv("FLOW_CACHE_PATH")
	Config.ServerAddress = os.Getenv("FLOW_SERVER_ADDRESS")
	Config.Name = os.Getenv("FLOW_NAME")

	if port, err := strconv.Atoi(os.Getenv("FLOW_LISTEN_PORT")); err == nil {
		Config.ListenPort = port
	}

	if debugBoolean, err := strconv.ParseBool(os.Getenv("FLOW_DEBUG")); err == nil {
		Config.Debug = debugBoolean
	}
}
