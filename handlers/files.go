package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"
	"retargetly-exercise/middleware"
	"retargetly-exercise/models"
)

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := models.APIResponse{}

	if r.URL.Path == "/files/list" || r.URL.Path == "/files/metrics" {
		//Check user token
		err := middleware.TokenValid(r)
		if err != nil {
			errMsg := "Unauthorized: " + err.Error()
			response.SendInfoMessage(w, errMsg, http.StatusUnauthorized)
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
		filesList, err := getFilesList(w, r, "./data")
		if err != nil {
			response.SendInfoMessage(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Content = filesList
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case "/files/metrics":
	}
}

func getFilesList(w http.ResponseWriter, r *http.Request, directory string) ([]models.FileListItem, error) {
	filesList, err := readFileNamesAndSize(directory, true)
	if err != nil {
		return []models.FileListItem{}, fmt.Errorf("Error while reading files at %s directory", directory)
	}
	return filesList, nil
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
			fileToAdd.Size = sizeFormating(info.Size())
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

func sizeFormating(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
