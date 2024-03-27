package transport

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/tahanasir/dicom-service/internal/image"
)

// Convert returns an HTTP handler function that converts DICOM files to PNG.
// It checks for a 'streaming' query parameter to decide between streaming conversion or not.
func Convert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")

		dicomFile, err := os.Open(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		info, err := dicomFile.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		streaming := r.URL.Query().Get("streaming")
		boolValue, err := strconv.ParseBool(streaming)

		// If file is greater than 1MB or 'streaming' is true and no error in parsing, perform streaming conversion.
		if info.Size() > 1<<20 || (err == nil && boolValue) {
			err = convertWithStreaming(dicomFile, info.Size(), filename)
		} else {
			err = convertWithoutStreaming(dicomFile, info.Size(), filename)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, fmt.Sprintf("%s.png", filename))
	}
}

// This is typically used for smaller files that can be processed in memory.
func convertWithoutStreaming(file io.Reader, size int64, filename string) error {
	dcm, err := dicom.Parse(file, size, nil)
	if err != nil {
		return err
	}

	// Find the PixelData element within the DICOM data.
	pixelDataElement, err := dcm.FindElementByTag(tag.PixelData)
	if err != nil {
		return err
	}

	// Extract pixel data information and write it to a PNG file.
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	image.WritePixelDataElement(pixelDataInfo, filename)
	return nil
}

// This method is preferred for large files to minimize memory usage.
func convertWithStreaming(file io.Reader, size int64, filename string) error {
	image.ParseWithStreaming(file, size, filename)
	return nil
}
