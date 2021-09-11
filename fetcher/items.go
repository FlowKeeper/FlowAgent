package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FlowKeeper/FlowAgent/v2/cache"
	"github.com/FlowKeeper/FlowAgent/v2/config"
	"github.com/FlowKeeper/FlowAgent/v2/webserver"
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
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

	logger.Info(loggingArea, fmt.Sprintf("Recieved %d items from server", len(agent.GetAllItems())))
	cache.RemoteAgent = agent
	webserver.ReadyToServe = true
	return nil
}
