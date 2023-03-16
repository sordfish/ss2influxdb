package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SSApiTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  string `json:"expires_in"`
	Scope        string `json:"scope"`
}

func GetAuthToken(user, pass string) SSApiTokenResponse {

	url := SSApiTokenEndpoint

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	d := SSApiTokenResponse{}
	jsonErr := json.Unmarshal(body, &d)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(d.AccessToken)

	return d
}
