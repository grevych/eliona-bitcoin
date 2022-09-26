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

package conf

import (
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

// Endpoint returns the configured API endpoint to get bitcoin rates data.
func Endpoint() string {
	return get("endpoint", "https://api.coindesk.com/v1/bpi/currentprice.json")
}

// Value returns the configuration string referenced by key. The configuration is stored in the init
// table bitcoin.configuration. This table should be configurable via the eliona frontend.
func get(name string, fallback string) string {
	valueChan := make(chan string)
	go func() {
		_ = db.Query(db.Pool(), "select value from bitcoin.configuration where name = $1", valueChan, name)
	}()
	value := <-valueChan
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Set sets the value of configuration
func Set(connection db.Connection, name string, value string) error {
	return db.Exec(connection, "insert into bitcoin.configuration (name, value) values ($1, $2) on conflict (name) do update set value = excluded.value", name, value)
}
