package scheduler

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/FlowKeeper/FlowAgent/v2/cache"
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const loggingArea = "Scheduler"

//StartScheduler gets called for every item seperatly and starts the related thread
func StartScheduler(ID primitive.ObjectID) {
	itemName := cache.CurrentItems[ID].Name
	logger.Info(loggingArea, "Starting scheduler for", itemName)

	for {
		item, found := cache.CurrentItems[ID]
		if !found {
			logger.Info(loggingArea, "Item", itemName, "disappeared. Scheduler now exiting")
			break
		}

		var result models.Result

		var cmd *exec.Cmd

		switch cache.RemoteAgent.OS {
		case models.Linux:
			cmd = exec.Command("bash", "-c", item.Command)
		case models.Windows:
			cmd = exec.Command("powershell", "-Command", item.Command)
		}

		rawOut, err := cmd.CombinedOutput()
		result.CapturedAt = time.Now()
		result.ItemID = ID
		result.Type = item.Returns

		if err != nil {
			//Couldn't execute
			result.Error = err.Error()
		} else {
			//Parse returned output from the command to the expected type (int or string)
			switch item.Returns {
			case models.Numeric:
				{
					//To avoid unnecessary problems with ParseInt
					saveOutput := strings.ReplaceAll(strings.ReplaceAll(string(rawOut), "\n", ""), "\r", "")
					value, err := strconv.ParseFloat(saveOutput, 64)
					if err != nil {
						result.Error = err.Error()
						break
					}

					result.ValueNumeric = value
					break
				}
			case models.Text:
				{
					result.ValueString = string(rawOut)
					break
				}
			default:
				{
					logger.Error(loggingArea, "Undefined item return type:", item.Returns)
					break
				}
			}
		}

		//logger.Debug("Result of item", item.Name, ":", strings.ReplaceAll(fmt.Sprint(result), "\n", " "))

		cache.AddResult(result)

		//Wait until checking again
		//Adding minimum delay
		if item.Interval < 1 {
			time.Sleep(time.Second)
		} else {
			time.Sleep(time.Duration(item.Interval) * time.Second)
		}
	}
}
