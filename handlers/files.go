package handlers

import (
	"net/http"
	"retargetly-exercise/models"
)

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := models.APIResponse{}

	if r.URL.Path == "/files/list" || r.URL.Path == "/files/metrics" {
		//Check user token
		if !isTokenValid(){
			response.SendInfoMessage(w, "Unauthorized, token is invalid", http.StatusUnauthorized)
			return
		}
		switch r.Method {
		case http.MethodGet:
			filesRoutesHandler(w, r, &response)
		default:
			// Response for other http Methods that are not GET
			response.SendInfoMessage(w, "Not implemented, try with GET method", http.StatusNotImplemented)
		}
	} else {
		response.SendInfoMessage(w, "Error: endpoint does not exist", http.StatusNotFound)
		return
	}
}

func filesRoutesHandler(w http.ResponseWriter, r *http.Request, response *models.APIResponse){
	switch r.URL.Path {
	case "/files/list":
	case "/files/metrics":
	}
}

func isTokenValid() bool{
	return false
}

func getFilesList(w http.ResponseWriter, r *http.Request) {

}

func getFilesMetrics(w http.ResponseWriter, r *http.Request) {

}
