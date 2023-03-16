package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
