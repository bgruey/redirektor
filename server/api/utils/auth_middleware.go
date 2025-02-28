package utils

import (
	"redirektor/server/repo"

	"net/http"
)

type AuthHandler struct {
	handler http.Handler
	psql    *repo.PostgresClient
}

func (mh *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("API-Key")
	if key == "" {
		RespondWithError(w, http.StatusBadRequest, "missing api key")
		return
	}

	storedKey, err := mh.psql.GetApiKey(key, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if storedKey == nil {
		RespondWithError(w, http.StatusBadRequest, "invalid api key")
		return
	}

	mh.handler.ServeHTTP(w, r)
}

func NewAuthHandler(handlerToWrap http.Handler, route string, mux *http.ServeMux) *AuthHandler {
	ret := &AuthHandler{
		handler: handlerToWrap,
		psql:    repo.NewPostgresClient(),
	}

	mux.Handle(route, ret)
	return ret
}
