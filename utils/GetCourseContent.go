package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"snek/types"
)

func GetCourseContent(host, apiKey, courseSlug string) (types.CheckoutResponse, error) {
	req, err := http.NewRequest("GET", host+"/api/checkout?slug="+courseSlug, nil)
	req.Header.Add("X-Api-Key", apiKey)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("An error occurred: Status " + resp.Status)
	}
	defer resp.Body.Close()
	var content types.CheckoutResponse
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		log.Fatal(err)
	}

	return content, nil
}
