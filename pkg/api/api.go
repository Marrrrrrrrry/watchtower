package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const tokenMissingMsg = "api token is empty or has not been set. exiting"

// DefaultListenAddress is the address the HTTP API listens on when no other
// address is configured.
const DefaultListenAddress = ":8080"

// API is the http server responsible for serving the HTTP API endpoints
type API struct {
	Token         string
	ListenAddress string
	hasHandlers   bool
}

// New is a factory function creating a new API instance
func New(token string, listenAddress string) *API {
	return &API{
		Token:         token,
		ListenAddress: NormalizeListenAddress(listenAddress),
		hasHandlers:   false,
	}
}

// NormalizeListenAddress turns a bare port such as "8080" into ":8080" as
// expected by http.ListenAndServe, defaults an empty address to
// DefaultListenAddress, and leaves "host:port" values untouched.
func NormalizeListenAddress(address string) string {
	switch {
	case address == "":
		return DefaultListenAddress
	case !strings.Contains(address, ":"):
		return ":" + address
	default:
		return address
	}
}

// RequireToken is wrapper around http.HandleFunc that checks token validity
func (api *API) RequireToken(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		want := fmt.Sprintf("Bearer %s", api.Token)
		if auth != want {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Debug("Valid token found.")
		fn(w, r)
	}
}

// RegisterFunc is a wrapper around http.HandleFunc that also sets the flag used to determine whether to launch the API
func (api *API) RegisterFunc(path string, fn http.HandlerFunc) {
	api.hasHandlers = true
	http.HandleFunc(path, api.RequireToken(fn))
}

// RegisterHandler is a wrapper around http.Handler that also sets the flag used to determine whether to launch the API
func (api *API) RegisterHandler(path string, handler http.Handler) {
	api.hasHandlers = true
	http.Handle(path, api.RequireToken(handler.ServeHTTP))
}

// Start the API and serve over HTTP. Requires an API Token to be set.
func (api *API) Start(block bool) error {

	if !api.hasHandlers {
		log.Debug("Watchtower HTTP API skipped.")
		return nil
	}

	if api.Token == "" {
		log.Fatal(tokenMissingMsg)
	}

	if block {
		api.runHTTPServer()
	} else {
		go func() {
			api.runHTTPServer()
		}()
	}
	return nil
}

func (api *API) runHTTPServer() {
	log.Fatal(http.ListenAndServe(api.ListenAddress, nil))
}
