package fetcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/cache"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/config"
)

func register() error {
	//If we already ran this function once
	if registered {
		return nil
	}

	logger.Info(loggingArea, "Trying to contact server for potential register")

	requestBody := struct {
		AgentUUID string
		AgentOS   string
		AgentPort string
		AgentName string
	}{
		AgentUUID: cache.Config.AgentUUID.String(),
		AgentOS:   runtime.GOOS,
		AgentPort: fmt.Sprint(config.Config.ListenPort),
		AgentName: config.Config.Name,
	}

	data, _ := json.Marshal(requestBody)

	resp, err := http.Post(fmt.Sprintf("%s://%s/api/v1/register", config.Config.ServerAdressParsed.Scheme, config.Config.ServerAdressParsed.Host), "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Error(loggingArea, "Couldn't register at server:", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		logger.Error(loggingArea, "Couldn't register at server:", string(respBytes))
		return errors.New("recieved invalid status code")
	}

	registered = true
	logger.Info(loggingArea, "My UUID is", cache.Config.AgentUUID)
	return nil
}
