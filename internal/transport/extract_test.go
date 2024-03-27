package transport

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		queryTag   string
		queryFile  string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid Tag and File",
			queryTag:   "(0010,0010)",
			queryFile:  "test-files/IM000001.dcm",
			wantStatus: http.StatusOK,
			wantBody:   "Value for tag (0010,0010): [NAYYAR^HARSH]",
		},
		{
			name:       "Invalid Tag",
			queryTag:   "(0010,0010",
			queryFile:  "test-files/IM000001.dcm",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid tag",
		},
		{
			name:       "File Not Found",
			queryTag:   "(0010,0010)",
			queryFile:  "test-files/IM000001.dc",
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Unable to parse DICOM file:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup request and recorder
			req, err := http.NewRequest("GET", fmt.Sprintf("/extract?tag=%s&filename=%s", tc.queryTag, tc.queryFile), nil)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()

			// Initialize handler and serve
			handler := http.HandlerFunc(Extract())
			handler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tc.wantStatus, rr.Code)

			// Assert body content
			assert.Contains(t, rr.Body.String(), tc.wantBody)
		})
	}
}
