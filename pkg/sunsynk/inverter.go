package sunsynk

import (
	"ssctl/pkg/utils"
)

var (
	SSAPIInverterEndpoint = "https://api.sunsynk.net/api/v1/inverter/"
)

type SSApiInverterGridRealtimeData struct {
	ETodayFrom        string  `json:"etodayFrom"`
	ETodayTo          string  `json:"etodayTo"`
	ETotalFrom        string  `json:"etotalFrom"`
	ETotalTo          string  `json:"etotalTo"`
	Fac               float64 `json:"fac"`
	LimiterTotalPower int     `json:"limiterTotalPower"`
	Pac               int     `json:"pac"`
	Pf                float64 `json:"pf"`
	Qac               int     `json:"qac"`
	Status            int     `json:"status"`
	Vip               []struct {
		Current string `json:"current"`
		Power   int    `json:"power"`
		Volt    string `json:"volt"`
	}
}

type SSApiInverterGridRealtimeDataResponse struct {
	Code    int                           `json:"code"`
	Message string                        `json:"msg"`
	Data    SSApiInverterGridRealtimeData `json:"data"`
	Success bool                          `json:"success"`
}

func GetInverterGridRealtimeData(inverterid, token string) ([]byte, error) {

	url := SSAPIInverterEndpoint + "grid/" + inverterid + "/realtime"

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := []byte{}
	respBody, err := utils.SendHTTPRequest("GET", url, headers, body, token)
	if err != nil {
		return respBody, err
	}

	return respBody, err

}

func GetInverterData(date, inverterid, column, token string) ([]byte, error) {

	url := SSAPIInverterEndpoint + "/energy/" + inverterid + "/input/day?lan=en&date=" + date + "&column=" + column

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := []byte{}
	respBody, err := utils.SendHTTPRequest("GET", url, headers, body, token)
	if err != nil {
		return respBody, err
	}

	return respBody, err

}

func GetCustomInverterData(date, edate, inverterid, params, token string) ([]byte, error) {

	url := SSAPIInverterEndpoint + inverterid + "/input/day?lan=en&date=" + date + "&edate=" + edate + "&params=" + params

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := []byte{}
	respBody, err := utils.SendHTTPRequest("GET", url, headers, body, token)
	if err != nil {
		return respBody, err
	}

	return respBody, err

}
