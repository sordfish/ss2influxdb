package sunsynk

import (
	"ssctl/pkg/utils"
)

var (
	SSApiPlantEndpoint    = "https://pv.inteless.com/api/v1/plant/"
	SSAPIInverterEndpoint = "https://pv.inteless.com/api/v1/inverter/"
)

type SSApiPlantData struct {
	Unit    string `json:"unit"`
	Records []struct {
		Time       string `json:"time"`
		Value      string `json:"value"`
		UpdateTime string `json:"updateTime"`
	} `json:"records"`
	Id        string `json:"id"`
	Label     string `json:"label"`
	GroupCode string `json:"groupCode"`
	Name      string `json:"name"`
}

type SSApiPlantDataResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Total int              `json:"total"`
		Infos []SSApiPlantData `json:"infos"`
	} `json:"data"`
	Success bool `json:"success"`
}

func GetPlantData(date, plantid, token string) ([]byte, error) {

	url := SSApiPlantEndpoint + "energy/" + plantid + "/day?lan=en&date=" + date

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
