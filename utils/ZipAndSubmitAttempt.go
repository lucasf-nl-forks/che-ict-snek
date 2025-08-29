package utils

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

type AttemptStatus struct {
	ExerciseID int    `json:"exerciseId"`
	Succeeded  bool   `json:"succeeded"`
	Status     string `json:"status"`
	Slug       string `json:"slug"`
	Output     string `json:"output"`
	Runtime    int    `json:"runtime"`
}

func ZipAndSubmitAttempt(directoryPath string, uploadEndpoint string, apiKey string) (AttemptStatus, error) {
	// Create a temporary zip file
	var result AttemptStatus

	tempDir, err := os.MkdirTemp("", "snek-attempt-*")
	tempFile := filepath.Join(tempDir, "upload.zip")
	fmt.Println(tempFile)
	if err != nil {
		return result, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               false,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	contents := []string{}
	elems, _ := os.ReadDir(directoryPath)
	for _, elem := range elems {
		contents = append(contents, elem.Name())
	}

	err = z.Archive(contents, tempFile)
	if err != nil {
		return result, err
	}

	// Reopen the zip file for reading
	file, err := os.Open(tempFile)
	if err != nil {
		return result, fmt.Errorf("failed to reopen zip file: %v", err)
	}
	defer file.Close()

	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the zip file to the form
	part, err := writer.CreateFormFile("file", "upload.zip")
	if err != nil {
		return result, fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the zip file to the multipart writer
	_, err = io.Copy(part, file)
	if err != nil {
		return result, fmt.Errorf("failed to copy file to multipart writer: %v", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return result, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Send request
	req, err := http.NewRequest("POST", uploadEndpoint, body)
	if err != nil {
		return result, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return result, fmt.Errorf("upload failed with status %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
