package redirect

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"redirektor/server/api/utils"
	"redirektor/server/model"
	"redirektor/server/qrcode"
	"redirektor/server/repo"
)

type LinkHandler struct {
	sync.Mutex
	psql *repo.PostgresClient
	host string
}

func NewLinkHandler() *LinkHandler {
	ret := new(LinkHandler)
	ret.psql = repo.NewPostgresClient()

	ret.host = strings.TrimRight(os.Getenv("HOST"), "/")

	return ret
}

func (rh *LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rh.post(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "invalid method")
	}

}

func (uh *LinkHandler) post(w http.ResponseWriter, r *http.Request) {
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

	redirect.OwnerKey = r.Header.Get(AuthHeaderKey)

	err = uh.psql.CreateRedirect(&redirect, nil)
	if err != nil {
		panic(err)
	}

	// double save to db because we minimize the number of characters
	// in the hash
	shortUrl := uh.host + "/" + redirect.Hash
	redirect.QRCode, err = qrcode.GenerateQRBytes(shortUrl)
	if err != nil {
		return
	}
	err = uh.psql.SaveRedirect(&redirect, nil)
	if err != nil {
		return
	}

	utils.RespondWithJSON(
		w, http.StatusCreated,
		map[string]any{
			"short_url":  shortUrl,
			"qrcode_b64": base64.StdEncoding.EncodeToString(redirect.QRCode),
		},
	)
}
