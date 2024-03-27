package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/suyashkumar/dicom"
)

// Extract returns an HTTP handler function that extracts and returns the value
// of a specified DICOM tag from a given DICOM file.
func Extract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")

		// Validate the format of the tag.
		if err := validateTag(tag); err != nil {
			http.Error(w, "Invalid tag", http.StatusBadRequest)
			return
		}

		file := r.URL.Query().Get("filename")
		data, err := dicom.ParseFile(file, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse DICOM file: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Search for the specified tag in the file's elements.
		for _, element := range data.Elements {
			if fmt.Sprintf("%v", element.Tag) == tag {
				fmt.Fprintf(w, "Value for tag %s: %v\n", tag, element.Value)
				return
			}
		}

		fmt.Fprintf(w, "Tag not found")
	}
}

func validateTag(tag string) error {
	// Check the tag format (length and parentheses)
	if len(tag) != 11 || tag[0] != '(' || tag[len(tag)-1] != ')' {
		return errors.New("invalid tag format")
	}

	// Remove parentheses and split by comma
	cleanTag := strings.Trim(tag, "()")
	cleanSlice := strings.Split(cleanTag, ",")

	// Ensure we have exactly two parts
	if len(cleanSlice) != 2 {
		return errors.New("invalid tag parts")
	}

	// Validate each part as hexadecimal
	_, err1 := strconv.ParseInt(cleanSlice[0], 16, 64)
	_, err2 := strconv.ParseInt(cleanSlice[1], 16, 64)

	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid hexadecimal values in tag. Part1: %s, Part2: %s", err1, err2)
	}

	return nil
}
