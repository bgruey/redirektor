package api

import (
	"fmt"
	"log"
	"net/http"
	"redirektor/server/api/redirect"
	"redirektor/server/api/utils"
)

type APIHandler struct {
	ServerMux *http.ServeMux
	redirect  *redirect.RedirectHandler
	update    *utils.AuthHandler
	port      string
}

func NewAPIHandler(port int) *APIHandler {
	ret := new(APIHandler)
	ret.ServerMux = http.NewServeMux()
	ret.port = fmt.Sprintf(":%d", port)

	ret.redirect = redirect.NewRedirectHandler(ret.ServerMux)

	ret.update = utils.NewAuthHandler(redirect.NewUpdatesHandler(), "/update", ret.ServerMux)

	// Pass server mux to register all paths for sub-handler
	return ret
}

func (a *APIHandler) Run() {
	fmt.Printf("Starting server, listening on port %s\n", a.port)
	log.Fatal(http.ListenAndServe(a.port, a.ServerMux))
}
