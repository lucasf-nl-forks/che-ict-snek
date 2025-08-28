package utils

import (
	"encoding/json"
	"net/http"
)

func CheckAttemptStatus(checkEndpoint, apiKey string) (AttemptStatus, error) {
	var result AttemptStatus
	req, err := http.NewRequest("GET", checkEndpoint, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("X-Api-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
