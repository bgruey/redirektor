package redirect

import (
	"net/http"

	"redirektor/server/api/utils"
	"redirektor/server/repo"
)

type AuthHandler struct {
	handler     http.Handler
	psql        *repo.PostgresClient
	requireRoot bool
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(AuthHeaderKey)
	if key == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing api key")
		return
	}

	storedKey, err := ah.psql.GetApiKey(key, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if storedKey == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid api key")
		return
	}

	if ah.requireRoot != storedKey.Root {
		var msg string
		if ah.requireRoot {
			msg = "root key required"
		} else {
			msg = "cannot use root key if not required"
		}
		utils.RespondWithError(w, http.StatusForbidden, msg)
	}

	ah.handler.ServeHTTP(w, r)
}

func NewAuthHandler(handlerToWrap http.Handler, route string, mux *http.ServeMux, requireRoot bool) *AuthHandler {
	ret := &AuthHandler{
		handler:     handlerToWrap,
		psql:        repo.NewPostgresClient(),
		requireRoot: requireRoot,
	}

	mux.Handle(route, ret)
	return ret
}
