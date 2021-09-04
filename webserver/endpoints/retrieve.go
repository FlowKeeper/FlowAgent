package endpoints

import (
	"net/http"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/cache"
)

func Retrieve(w http.ResponseWriter, r *http.Request) {
	results, err := cache.RetrieveCache()
	if err != nil {
		httpResponse.InternalError(w, r, err)
		return
	}

	httpResponse.SuccessWithPayload(w, "OK", results)
}
