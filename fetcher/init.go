package fetcher

import (
	"runtime"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/cache"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/scheduler"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Init runs starts periodically fetching all needed items from the server
func Init() {
	cache.CurrentItems = make(map[primitive.ObjectID]models.Item)
	for {
		if err := fetch(); err == nil {
			syncConfig()
		}

		time.Sleep(60 * time.Second)
	}
}

func syncConfig() {
	newItems := make(map[primitive.ObjectID]models.Item, 0)
	unthreadedItems := make([]models.Item, 0)
	ignoredDueToOS := 0
	ignoredDueToDisabled := 0

	for _, k := range cache.RemoteAgent.ItemsResolved {
		agentos, _ := models.AgentosFromString(runtime.GOOS)
		if k.CheckOn != agentos {
			ignoredDueToOS++
			continue
		}
		/*
			if !k.Enabled {
				ignoredDueToDisabled++
				continue
			}
		*/

		//Copy all new items to new map
		newItems[k.ID] = k

		//Find out if the item was schedules before. If not we need to start a thread for it
		if _, found := cache.CurrentItems[k.ID]; !found {
			unthreadedItems = append(unthreadedItems, k)
			logger.Info(loggingArea, "We need to start a new thread for", k.Name)
		}

	}

	//Set new map to items map as these are now the current config
	cache.CurrentItems = newItems

	//Start a new thread for all new items
	for _, k := range unthreadedItems {
		go scheduler.StartScheduler(k.ID)
	}

	logger.Info(loggingArea, "Ignored", ignoredDueToOS, "Items due to agent OS")
	logger.Info(loggingArea, "Ignored", ignoredDueToDisabled, "Items due to being deactivated")
}
