package transport

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Upload returns an HTTP handler function for uploading DICOM files.
// It parses multipart/form-data requests to retrieve the uploaded file.
func Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit the size of the form to 10MB
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, fmt.Sprintf("File is larger than 10MB: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Retrieve the file from the form data.
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid file: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Create a new file in the images/dicom directory with the uploaded file's name.
		newFilename := fmt.Sprintf("%s.dcm", fileHeader.Filename)
		newFile, err := os.Create(fmt.Sprintf("images/dicom/%s", newFilename))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to create file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to read file: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if _, err = newFile.Write(fileBytes); err != nil {
			http.Error(w, fmt.Sprintf("Unable to write file: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Successfully uploaded", newFile.Name())
	}
}
