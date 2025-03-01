package redirect

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"redirektor/server/api/utils"
	"redirektor/server/model"
	"redirektor/server/repo"
)

type DeleteRequest struct {
	Key       string `json:"api_key"`
	DeletedAt int64  `json:"deleted_at"`
}

type ApiKeyHandler struct {
	sync.Mutex
	psql       *repo.PostgresClient
	rootApiKey *model.ApiKey
}

func NewApiKeyHandler() *ApiKeyHandler {
	ret := new(ApiKeyHandler)
	ret.psql = repo.NewPostgresClient()

	key, err := ret.psql.GetRootKey(nil)
	if err != nil {
		panic(err)
	}

	ret.rootApiKey = key
	log.Printf("(Root) %s: %s\n", AuthHeaderKey, ret.rootApiKey.Key)

	return ret
}

func (akh *ApiKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		akh.post(w)
	case "DELETE":
		akh.delete(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}

}

func (akh *ApiKeyHandler) delete(w http.ResponseWriter, r *http.Request) {
	defer func() {
		akh.Unlock()
		r.Body.Close()
	}()

	akh.Lock()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to read body, maybe retry?")
		return
	}

	deleteRequest := &DeleteRequest{}
	err = json.Unmarshal(body, deleteRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "could not parse request")
		return
	}
	if deleteRequest.Key == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "no key supplied")
		return
	}

	err = akh.psql.DeleteKey(deleteRequest.Key, deleteRequest.DeletedAt, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to delete key")
		return
	}

	utils.RespondWithJSON(
		w, http.StatusOK,
		map[string]any{
			"api-key":    deleteRequest.Key,
			"deleted_at": deleteRequest.DeletedAt,
		},
	)
}

func (akh *ApiKeyHandler) post(w http.ResponseWriter) {
	defer akh.Unlock()
	akh.Lock()

	key := model.NewApiKey()
	err := akh.psql.CreateApiKey(key, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "sorry")
		return
	}

	utils.RespondWithJSON(
		w, http.StatusCreated,
		map[string]string{"api_key": key.Key},
	)

}
