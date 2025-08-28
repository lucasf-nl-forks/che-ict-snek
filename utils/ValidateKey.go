package utils

import (
	"fmt"
	"net/http"
)

func ValidateKey(domain, key string) error {
	req, err := http.NewRequest("GET", domain+"/api/auth/validate", nil)
	req.Header.Add("X-Api-Key", key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Login invalid: %s", res.Status)
	}
	return nil
}
