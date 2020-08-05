package models

import (
	"encoding/json"
	"net/http"
)

// #region API Responses structs
type APIResponse struct {
	Content interface{} `json:"response"`
}

type InfoResponseItem struct {
	Message string `json:"message"`
}

type LoginResponseItem struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type FileListResponse struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func (res *APIResponse) SendInfoMessage(w http.ResponseWriter, message string, status int) {
	infoMessage := InfoResponseItem{
		Message: message,
	}
	(*res).Content = infoMessage
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*res)
}

//#endregion
