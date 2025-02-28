package redirect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"redirektor/server/api/utils"
	"redirektor/server/model"
	"redirektor/server/repo"
	"strings"
	"sync"
)

type UpdatesHandler struct {
	sync.Mutex
	psql     *repo.PostgresClient
	Password string
	host     string
}

func NewUpdatesHandler() *UpdatesHandler {
	ret := new(UpdatesHandler)
	ret.psql = repo.NewPostgresClient()

	ret.Password = os.Getenv("API_PASSWORD")
	ret.host = strings.TrimRight(os.Getenv("HOST"), "/")

	key := model.NewApiKey()
	count, err := ret.psql.CountApiKeys(nil)
	if err != nil {
		panic(err)
	}
	if count < 1 {
		ret.psql.CreateApiKey(key, nil)
	} else {
		key, err = ret.psql.GetSingleApiKey(nil)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("API-Key: %s\n", key.Key)

	return ret
}

func (rh *UpdatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)

	switch r.Method {
	case "POST":
		rh.post(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}

}

func (uh *UpdatesHandler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to read body, maybe retry?")
		return
	}

	defer uh.Unlock()
	uh.Lock()

	var redirect model.Redirect
	err = json.Unmarshal(body, &redirect)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json data")
	}

	err = uh.psql.CreateRedirect(&redirect, nil)
	if err != nil {
		panic(err)
	}

	utils.RespondWithJSON(
		w, http.StatusCreated,
		map[string]string{"short-url": uh.host + "/" + redirect.Hash},
	)
}
