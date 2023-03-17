package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	SSApiTokenEndpoint = "https://pv.inteless.com/oauth/token"
)

type SSApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Success bool   `json:"success"`
}

type SSApiTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  string `json:"expires_in"`
	Scope        string `json:"scope"`
}

func GetAuthToken(user, pass string) SSApiTokenResponse {

	fmt.Println("Getting token")

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

func livez(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func healthz(w http.ResponseWriter, req *http.Request, c mqtt.Client) {
	if c.IsConnected() {
		fmt.Fprintf(w, "OK\n")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Not connected to mqtt"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error with JSON marshal - Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

}

func main() {

	brokerPtr := flag.String("broker", "tcp://127.0.0.1:1883", "MQTT broker address")
	flag.Parse()

	ss_user := os.Getenv("SS_USER")
	ss_pass := os.Getenv("SS_PASS")

	if ss_pass == "" || ss_user == "" {
		log.Fatalf("No creds defined")
	}

	ss_auth_token := GetAuthToken(ss_user, ss_pass)

	fmt.Println(ss_auth_token.AccessToken)
	fmt.Println(ss_auth_token.RefreshToken)
	fmt.Println(ss_auth_token.Scope)
	fmt.Println(ss_auth_token.TokenExpiry)
	fmt.Println(ss_auth_token.TokenType)

	// Connect to the MQTT server
	opts := mqtt.NewClientOptions().AddBroker(*brokerPtr)
	opts.SetClientID("ss2mqtt")
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error with MQTT - Err: %s", token.Error())
	}

	http.HandleFunc("/livez", livez)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthz(w, r, c)
	})

	http.ListenAndServe(":34567", nil)

	c.Disconnect(250)

}
