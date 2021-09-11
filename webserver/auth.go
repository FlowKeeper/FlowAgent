package webserver

import (
	"net/http"

	"github.com/FlowKeeper/FlowAgent/v2/cache"
	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

const NotAllowed = "You are not allowed to access this agent"

var ReadyToServer = false

//authorizationMiddleware should check the "ScraperUUID" header and determine if the client is allowed to send http requests to this agent
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ReadyToServer {
			httpResponse.UserError(w, 406, "Please wait until initialization is completed")
			logger.Warning(loggingArea, "Someone tried to access this agent pre-initializtaion!")
			return
		}

		scraperuuidString := r.Header.Get("ScraperUUID")
		if stringHelper.IsEmpty(scraperuuidString) {
			httpResponse.UserError(w, 401, NotAllowed)
			logger.Warning(loggingArea, "Someone tried to access this agent without providing a ScraperUUID")
			return
		}
		scraperUUID, err := uuid.Parse(scraperuuidString)
		if err != nil {
			httpResponse.UserError(w, 401, NotAllowed)
			logger.Warning(loggingArea, "Someone tried to access this agent with an invalid ScraperUUID header")
			return
		}

		if cache.RemoteAgent.Scraper.UUID != scraperUUID {
			httpResponse.UserError(w, 401, NotAllowed)
			logger.Warning(loggingArea, "Someone tried to access this agent with the wrong ScraperUUID")
			return
		}
		next.ServeHTTP(w, r)
	})
}
