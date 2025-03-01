package api

import (
	"fmt"
	"log"
	"net/http"
	"redirektor/server/api/redirect"
)

type APIHandler struct {
	ServerMux *http.ServeMux
	redirect  *redirect.RedirectHandler
	link      *redirect.AuthHandler
	key       *redirect.AuthHandler
	port      string
}

func NewAPIHandler(port int) *APIHandler {
	ret := new(APIHandler)
	ret.ServerMux = http.NewServeMux()
	ret.port = fmt.Sprintf(":%d", port)

	// general usage
	ret.redirect = redirect.NewRedirectHandler(ret.ServerMux)
	// add links for redirects
	ret.link = redirect.NewAuthHandler(redirect.NewLinkHandler(), "/link", ret.ServerMux, false)
	// add/delete api keys
	ret.key = redirect.NewAuthHandler(redirect.NewApiKeyHandler(), "/key", ret.ServerMux, true)

	// Pass server mux to register all paths for sub-handler
	return ret
}

func (a *APIHandler) Run() {
	fmt.Printf("Starting server, listening on port %s\n", a.port)
	log.Fatal(http.ListenAndServe(a.port, a.ServerMux))
}
