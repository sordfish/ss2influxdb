package sunsynk

import (
	"ssctl/pkg/utils"
	"time"
)

var (
	SSApiPlantEndpoint = "https://api.sunsynk.net/api/v1/plant/"
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

type SSApiPlantInverterData struct {
	Sn           string `json:"sn"`
	Alias        string `json:"alias"`
	Gsn          string `json:"gsn"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
	CommTypeName string `json:"commTypeName"`
	CustCode     int    `json:"custCode"`
	Version      struct {
		MasterVer string `json:"masterVer"`
		SoftVer   string `json:"softVer"`
		HardVer   string `json:"hardVer"`
		HmiVer    string `json:"hmiVer"`
		BmsVer    string `json:"bmsVer"`
	} `json:"version"`
	Model     string    `json:"model"`
	EquipMode any       `json:"equipMode"`
	Pac       int       `json:"pac"`
	Etoday    float64   `json:"etoday"`
	Etotal    float64   `json:"etotal"`
	UpdateAt  time.Time `json:"updateAt"`
	Opened    int       `json:"opened"`
	Plant     struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Type      int    `json:"type"`
		Master    any    `json:"master"`
		Installer any    `json:"installer"`
		Email     any    `json:"email"`
		Phone     any    `json:"phone"`
	} `json:"plant"`
	GatewayVO struct {
		Gsn    string `json:"gsn"`
		Status int    `json:"status"`
	} `json:"gatewayVO"`
	SunsynkEquip       bool   `json:"sunsynkEquip"`
	ProtocolIdentifier string `json:"protocolIdentifier"`
	EquipType          int    `json:"equipType"`
	RatePower          int    `json:"ratePower"`
}

type SSApiPlantInverterDataResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		PageSize   int                      `json:"pageSize"`
		PageNumber int                      `json:"pageNumber"`
		Total      int                      `json:"total"`
		Infos      []SSApiPlantInverterData `json:"infos"`
	} `json:"data"`
	Success bool `json:"success"`
}

func GetInverterId(plantid, token string) ([]byte, error) {

	url := SSApiPlantEndpoint + plantid + "/inverters?page=1&limit=10&status=-1&type=-2"

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
