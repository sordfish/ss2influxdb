package sunsynk

import (
	"ssctl/pkg/utils"
)

var (
	SSApiListPlantsEndpoint = "https://api.sunsynk.net/api/v1/plants?page=1&limit=10"
)

type SSAuthToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Scope        string
	TokenExpiry  string
	Timestamp    string
}

type SSApiUserPlant struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SSApiUserPlantsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Total int              `json:"total"`
		Infos []SSApiUserPlant `json:"infos"`
	} `json:"data"`
	Success bool `json:"success"`
}

func GetUserData(url, token string) ([]byte, error) {

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
