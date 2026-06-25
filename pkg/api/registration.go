package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/joohoi/acme-dns/pkg/acmedns"
)

// DeleteRegistrationResponse is the JSON body for a successful DELETE /registration.
type DeleteRegistrationResponse struct {
	Message   string `json:"message"`
	Subdomain string `json:"subdomain"`
}

func (a *AcmednsAPI) webRegistrationDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, ok := r.Context().Value(ACMETxtKey).(acmedns.ACMETxt)
	if !ok {
		a.Logger.Errorw("Context error",
			"error", "context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonError("internal_error"))
		return
	}

	err := a.DB.Unregister(user.Username)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(jsonError("not_found"))
			return
		}
		a.Logger.Errorw("Error while trying to delete registration",
			"error", err.Error(),
			"user", user.Username.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonError("db_error"))
		return
	}

	a.Logger.Infow("Deleted registration",
		"user", user.Username.String(),
		"subdomain", user.Subdomain)

	resp := DeleteRegistrationResponse{
		Message:   "deleted",
		Subdomain: user.Subdomain,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonError("json_error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
