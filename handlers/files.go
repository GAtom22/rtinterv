package handlers

import (
	"encoding/json"
	"retargetly-exercise/models"
	"net/http"
)

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := models.APIResponse{}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/files/list" || r.URL.Path == "/files/metrics" {
			//Check user token

		switch r.URL.Path{
		case "/files/list":
		case "/files/metrics":
		}

		}else {
			response.SendInfoMessage(w, "Error: endpoint does not exist", http.StatusNotFound)
			return
		}


		//if OK generate and send auth token with expiration date
		// Generate token
		token, expirationDate, err := createToken(user, 10)
		if err != nil {
			response.SendInfoMessage(w, "Error while generating token", http.StatusInternalServerError)
			return
		}

		// send data to client
		loginData := models.LoginResponseItem{
			Token:   token,
			Expires: expirationDate,
		}
		response.Content = loginData
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		// Response for other http Methods that are not POST
		response := models.APIResponse{}
		response.SendInfoMessage(w, "Not implemented, try with GET method", http.StatusNotImplemented)
	}
}

func getFilesList(w http.ResponseWriter, r *http.Request) {

}

func getFilesMetrics(w http.ResponseWriter, r *http.Request) {

}
