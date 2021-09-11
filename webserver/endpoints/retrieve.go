package endpoints

import (
	"net/http"

	"github.com/FlowKeeper/FlowAgent/v2/cache"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
)

//Retrieve is called by the server to fetch the cached results
func Retrieve(w http.ResponseWriter, r *http.Request) {
	results, err := cache.RetrieveCache()
	if err != nil {
		httpResponse.InternalError(w, r, err)
		return
	}

	httpResponse.SuccessWithPayload(w, "OK", results)
}
