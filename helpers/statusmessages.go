package helpers

import (
	"encoding/json"
	"net/http"
	"retargetly-exercise/models"
	"time"
)

//SendStatusMessage for the file metrics section
func SendStatusMessage(response *models.APIResponse, w http.ResponseWriter, status string) int64 {
	reqTime := time.Now().Unix()
	statusMsg := models.FileMetricsStatusMessage{
		Status:    status,
		StartTime: FormatDate(reqTime),
	}
	(*response).Content = statusMsg
	json.NewEncoder(w).Encode(*response)
	return reqTime
}

//SendFailedStatusMessage sends status: failed message for the file metrics section
func SendFailedStatusMessage(response *models.APIResponse, w http.ResponseWriter, err error) {
	errMessage := models.FileMetricsFailedMessage{
		Status:  "failed",
		Message: err.Error(),
	}
	response.Content = errMessage
	json.NewEncoder(w).Encode(response)
}
