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
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/webserver"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
)

const loggingArea = "Fetcher"

//Registered is set to true after register() is run once, so we dont run it multiple times (would be pointless)
var registered = false

func fetch() error {
	//Ensure that we are registered
	if err := register(); err != nil {
		return err
	}

	logger.Info(loggingArea, "Fetching current config from server")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/api/v1/config", config.Config.ServerAdressParsed.Scheme, config.Config.ServerAdressParsed.Host), nil)
	req.Header.Add("AgentUUID", cache.Config.AgentUUID.String())

	if err != nil {
		logger.Fatal(loggingArea, "Couldn't construct server request:", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(loggingArea, "Couldn't send server request:", err)
		return err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error(loggingArea, "Couldn't retrieve items from server. Got:", string(bodyBytes))
		return errors.New("recieved invalid status code")
	}

	var parsedResponse struct {
		Status  string
		Payload models.Agent
	}

	if err := json.Unmarshal(bodyBytes, &parsedResponse); err != nil {
		logger.Error(loggingArea, "Couldn't deode response from server:", err, "Got:", string(bodyBytes))
		return err
	}

	agent := parsedResponse.Payload

	logger.Info(loggingArea, fmt.Sprintf("Recieved %d items from server", len(agent.ItemsResolved)))
	cache.RemoteAgent = agent
	webserver.ReadyToServer = true
	return nil
}
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
	}{
		AgentUUID: cache.Config.AgentUUID.String(),
		AgentOS:   runtime.GOOS,
		AgentPort: fmt.Sprint(config.Config.ListenPort),
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
