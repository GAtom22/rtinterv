package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"retargetly-exercise/dataprocessing"
	"retargetly-exercise/helpers"
	"retargetly-exercise/middleware"
	m "retargetly-exercise/models"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

//FilesHandler handles the request for files/list and files/metrics
func FilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := m.APIResponse{}

	if r.URL.Path == "/files/list" || r.URL.Path == "/files/metrics" {
		//Check user token
		err := middleware.TokenValid(r)
		if err != nil {
			errMsg := "Unauthorized: " + err.Error()
			response.SendInfoMessage(w, errMsg, http.StatusForbidden)
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

func filesRoutesHandler(w http.ResponseWriter, r *http.Request, response *m.APIResponse) {
	filesDirectory := "./data"
	switch r.URL.Path {

	case "/files/list":
		filesList, status, err := getFilesList(w, r, filesDirectory)
		if err != nil {
			response.SendInfoMessage(w, err.Error(), status)
			return
		}
		response.Content = filesList
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)

	case "/files/metrics":
		//Create a channel to receive the metrics data
		metricsCh := make(chan *m.FileMetricsSucessMessage, 10)
		//Check if the request is ok (filename parameter) - send failed message if something is wrong
		fileLocation, err := getFileNameFromReq(r, filesDirectory)
		if err != nil {
			helpers.SendFailedStatusMessage(response, w, err)
			return
		}

		// Record the time of the request and send status started message
		reqTime := helpers.SendStatusMessage(response, w, "started")

		//flush to send updates of the request status to the client
		f, ok := w.(http.Flusher)
		if ok {
			f.Flush()
		} else {
			fmt.Println("http Flush not working")
		}

		//Update status to processing and send it to client
		helpers.SendStatusMessage(response, w, "processing")
		f.Flush()

		//Create go routine to perform the data processing - send the data thru the channel when done
		go func() {
			//Do the data processing and get the metrics
			fileMetricsResponse, _, err := getFilesMetrics(w, r, fileLocation, reqTime)
			if err != nil {
				helpers.SendFailedStatusMessage(response, w, err)
				return
			}
			metricsCh <- &fileMetricsResponse
		}()

		// Return a successfull response with the data
		response.Content = <-metricsCh
		json.NewEncoder(w).Encode(response)
	}
}

func getFilesList(w http.ResponseWriter, r *http.Request, directory string) ([]m.FileListItem, int, error) {
	// Check if humanreadable is requested to get the file sizes in friendly units
	keys, ok := r.URL.Query()["humanreadable"]
	humanReadable := false
	if ok && keys[0] == "true" {
		humanReadable = true
	}
	// Go to the directory and get the files list
	filesList, err := readFileNamesAndSize(directory, humanReadable)
	if err != nil {
		return []m.FileListItem{}, http.StatusInternalServerError, fmt.Errorf("Error while fetching files at %s directory", directory)
	}
	if len(filesList) == 0 {
		return []m.FileListItem{}, http.StatusOK, fmt.Errorf("No files found at %s directory", directory)
	}
	return filesList, http.StatusOK, nil
}

func getFileNameFromReq(r *http.Request, directory string) (string, error) {
	filesList, err := readFileNamesAndSize(directory, false)

	if err != nil {
		return "", fmt.Errorf("Error while fetching files at %s directory", directory)
	}
	// get the filename and check if it is in the data directory
	fileName, ok := r.URL.Query()["filename"]

	if !ok || !isInFileList(fileName[0], filesList) {
		errorMsg := fmt.Errorf(`Error, "filename" parameter not present in the request`)

		if len(fileName) > 0 {
			errorMsg = fmt.Errorf("Error, file %s not found in directory %s", fileName[0], directory)
		}

		return "", errorMsg
	}
	fileLocation := fmt.Sprintf("%s/%s", directory, fileName[0])
	return fileLocation, nil
}

func getFilesMetrics(w http.ResponseWriter, r *http.Request, fileLocation string, reqTime int64) (m.FileMetricsSucessMessage, int, error) {
	response := m.FileMetricsSucessMessage{}
	response.StartTime = helpers.FormatDate(reqTime)

	endTime, fileMetrics, err := dataprocessing.GetFileMetrics(fileLocation)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}
	response.Status = "ready"
	response.Metrics = fileMetrics
	response.EndTime = helpers.FormatDate(endTime)
	return response, http.StatusOK, nil
}

func readFileNamesAndSize(root string, humanReadable bool) ([]m.FileListItem, error) {
	var files []m.FileListItem
	// Read files at "root" directory
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Skip if it is a directory
		if info.IsDir() {
			return nil
		}
		fileToAdd := m.FileListItem{
			Name: info.Name(),
		}
		// Format file size if requested
		if humanReadable {
			fileToAdd.Size = helpers.FileSizeFormating(info.Size())
		} else {
			fileToAdd.Size = info.Size()
		}
		files = append(files, fileToAdd)
		return nil
	})

	if err != nil {
		return []m.FileListItem{}, err
	}
	return files, nil
}

func isInFileList(fileName string, filesList []m.FileListItem) bool {
	for _, file := range filesList {
		if strings.TrimSpace(fileName) == file.Name {
			return true
		}
	}
	return false
}
