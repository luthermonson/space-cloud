package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleSetUserManagement returns the handler to get the project config and validate the user via a REST endpoint
func HandleSetUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		provider := vars["id"]

		// Load the body of the request
		value := new(config.AuthStub)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := adminMan.IsTokenValid(token, "auth-provider", "modify", map[string]string{"project": projectID, "id": provider})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = value
		if status, err := syncMan.SetUserManagement(ctx, projectID, provider, value, reqParams); err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		// Give a positive acknowledgement
		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetUserManagement returns handler to get auth
func HandleGetUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		providerID := "*"
		providerQuery, exists := r.URL.Query()["id"]
		if exists {
			providerID = providerQuery[0]
		}

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(token, "auth-provider", "modify", map[string]string{"project": projectID, "id": providerID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		status, providers, err := syncMan.GetUserManagement(ctx, projectID, providerID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}
		_ = utils.SendResponse(w, status, model.Response{Result: providers})
	}
}
