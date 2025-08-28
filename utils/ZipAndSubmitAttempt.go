package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	tempFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		return result, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write directory contents to zip file
	zw := zip.NewWriter(tempFile)
	err = filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %v", err)
		}

		relPath, _ := filepath.Rel(directoryPath, path)
		if relPath == "." {
			return nil
		}

		info, _ := d.Info()
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %v", err)
		}
		header.Method = zip.Deflate
		header.Name = relPath

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %v", err)
		}

		if !d.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("failed to write file to zip: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	err = zw.Close()
	if err != nil {
		return result, fmt.Errorf("failed to close zip writer: %v", err)
	}

	// Reopen the zip file for reading
	file, err := os.Open(tempFile.Name())
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
