package handlers

import (
	"retargetly-exercise/helpers"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"retargetly-exercise/middleware"
	"retargetly-exercise/models"
)

//FilesHandler handles the request for files/list and files/metrics
func FilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := models.APIResponse{}

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

func filesRoutesHandler(w http.ResponseWriter, r *http.Request, response *models.APIResponse) {
	switch r.URL.Path {
	case "/files/list":
		filesList, status, err := getFilesList(w, r, "./data")
		if err != nil {
			response.SendInfoMessage(w, err.Error(), status)
			return
		}
		response.Content = filesList
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case "/files/metrics":
	}
}

func getFilesList(w http.ResponseWriter, r *http.Request, directory string) ([]models.FileListItem, int, error) {
	// Check if humanreadable is requested to get the file sizes in friendly units
	keys, ok := r.URL.Query()["humanreadable"]
	humanReadable := false
	if ok && keys[0] == "true" {
		humanReadable = true
	}
	// Go to the directory and get the files list
	filesList, err := readFileNamesAndSize(directory, humanReadable)
	if err != nil {
		return []models.FileListItem{}, http.StatusInternalServerError, fmt.Errorf("Error while reading files at %s directory", directory)
	}
	if len(filesList) == 0 {
		return []models.FileListItem{}, http.StatusOK, fmt.Errorf("No files found at %s directory", directory)
	}
	return filesList, http.StatusOK, nil
}

func getFilesMetrics(w http.ResponseWriter, r *http.Request) {

}

func readFileNamesAndSize(root string, humanReadable bool) ([]models.FileListItem, error) {
	var files []models.FileListItem
	// Read files at "root" directory
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Skip if it is a directory
		if info.IsDir() {
			return nil
		}
		fileToAdd := models.FileListItem{
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
		return []models.FileListItem{}, err
	}
	return files, nil
}
