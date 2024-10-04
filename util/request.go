package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func MakeAuthRequest(method, url string, body io.Reader) (*http.Response, error) {
	// Get the authorization token from viper
	token := viper.GetString("token")
	if token == "" {
		return nil, fmt.Errorf("authorization token is missing")
	}

	// Create the request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add the headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	// Create the HTTP client
	client := &http.Client{}

	// Make the request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return resp, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized request")
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden request")
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func MakeUnAuthRequest(method, url string, body io.Reader) (*http.Response, error) {
	// Create the request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add the headers
	req.Header.Add("Content-Type", "application/json")

	// Create the HTTP client
	client := &http.Client{}

	// Make the request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return resp, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized request")
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden request")
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

}

func MakeAuthRequestAndParseResponse(method, url string, jsonBody interface{}, response interface{}) error {

	var body io.Reader

	if jsonBody != nil {

		jsonData, err := json.Marshal(jsonBody)

		if err != nil {
			return err
		}

		body = bytes.NewBuffer(jsonData)
	}

	resp, err := MakeAuthRequest(method, url, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, response)
}

func MakeUnAuthRequestAndParseResponse(method, url string, jsonBody interface{}, response interface{}) error {

	var body io.Reader

	if jsonBody != nil {

		jsonData, err := json.Marshal(jsonBody)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(jsonData)
	}

	resp, err := MakeUnAuthRequest(method, url, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, response)
}
