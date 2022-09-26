//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"bitcoin/apiserver"
	"bitcoin/apiservice"
	"net/http"

	"github.com/eliona-smart-building-assistant/go-utils/common"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// listenApiRequests starts an API server and listen for API requests
func listenApiRequests() {
	port := ":" + common.Getenv("BITCOIN_INTEGRATION_PORT", "3001")
	err := http.ListenAndServe(port, apiserver.NewRouter(
		apiserver.NewApiController(apiservice.NewApiService()),
	))
	log.Fatal("Bitcoin", "Error in API Server: %v", err)
}
