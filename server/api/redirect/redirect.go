package redirect

import (
	"net/http"
	"sync"

	"redirektor/server/api/utils"
	"redirektor/server/repo"
)

type RedirectHandler struct {
	sync.Mutex
	psql *repo.PostgresClient
}

func NewRedirectHandler(mutex *http.ServeMux) *RedirectHandler {
	ret := new(RedirectHandler)
	ret.psql = repo.NewPostgresClient()

	mutex.Handle("/", ret)

	return ret
}

func (rh *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rh.get(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}

}

func (rh *RedirectHandler) get(w http.ResponseWriter, r *http.Request) {
	defer rh.Unlock()
	rh.Lock()

	hash, err := utils.HashFromUrl("/", r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid address")
		return
	}

	linkAddress, err := rh.psql.GetIncrementRedirectByHash(hash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "oops")
		return
	}
	if linkAddress == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "blank address"})
		return
	}

	http.Redirect(w, r, linkAddress, http.StatusPermanentRedirect)
}
