package models

import (
	"encoding/json"
	"net/http"
)

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

type FileListItem struct {
	Name string      `json:"name"`
	Size interface{} `json:"size"`
}

type FileMetricsStatusMessage struct{
	Status string `json:"status"`
	StartTime string `json:"started"`
}

type FileMetricsFailedMessage struct{
	Status string `json:"status"`
	Message string `json:"message"`
}

type FileMetricsSucessMessage struct{
	Status string `json:"status"`
	StartTime string `json:"started"`
	EndTime string `json:"finished"`
	Metrics []Segment `json:"metrics"`
}

type Segment struct{
	SegmentID int `json:"segmentId"`
	Uniques []Unique `json:"uniques"`
}

type Unique struct{
	Country string `json:"country"`
	Count int `json:"count"`
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
