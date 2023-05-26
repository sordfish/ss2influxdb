package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func SendHTTPRequest(method string, url string, headers map[string]string, body []byte, authToken string) ([]byte, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Add headers to the HTTP request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add the authentication token to the request header
	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the HTTP response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check the HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	debugEnabled := os.Getenv("SS_DEBUG")

	if debugEnabled == "TRUE" {
		log.Debug(string(respBody))
	}

	return respBody, nil
}

func Upload2influxdb(data string) {

	InfluxdbUrl := os.Getenv("INFLUXDB_URL")

	if InfluxdbUrl == "" {
		log.Fatal("InfluxDB not url set")
	}

	url := InfluxdbUrl + "/write"

	headers := map[string]string{}
	body := []byte(data)
	token := ""
	respBody, err := SendHTTPRequest("POST", url, headers, body, token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(respBody)

}
