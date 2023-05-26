package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"ssctl/pkg/kube"
	"ssctl/pkg/sunsynk"

	log "github.com/sirupsen/logrus"
)

type LineFormat struct {
	Value     float64
	Name      string
	Unit      string
	PlantId   int
	Timestamp int64
}

func GetToken(k8s bool) string {

	var SunsynkToken string

	if !k8s {

		SunsynkToken = os.Getenv("SS_TOKEN")

		if SunsynkToken == "" {
			log.Fatal("No token found in env 1")
		}

		return SunsynkToken

	} else {

		clientset, err := kube.Login()
		if err != nil {
			log.Fatal(err)
		}

		result, err := kube.GetK8sSecret(clientset, "sunsynk-token", "sunsynk")
		if err != nil {
			log.Fatal(err)
		}

		token, ok := result.Data["token"]
		if !ok {
			log.Fatal("token not found in secret data")
		}

		SunsynkToken = string(token)

		return SunsynkToken
	}

}

func GetPlantIDs(k8s bool) string {

	var SunsynkToken, SunsynkPlantId string

	if !k8s {

		SunsynkPlantId = os.Getenv("SS_PLANT_ID")
		SunsynkToken = os.Getenv("SS_TOKEN")

		if SunsynkPlantId == "" {
			log.Fatal("No plant ID found in env")
		}

		if SunsynkToken == "" {
			log.Fatal("No token found in env 2")
		}

	} else {

		clientset, err := kube.Login()
		if err != nil {
			log.Fatal(err)
		}

		result, err := kube.GetK8sSecret(clientset, "sunsynk-user-plants", "sunsynk")
		if err != nil {
			log.Fatal(err)
		}

		plantdata, ok := result.Data["plants.json"]
		if !ok {
			log.Fatal("plants.json not found in secret data")
		}

		var UserPlantsStruct []sunsynk.SSApiUserPlant

		err = json.Unmarshal(plantdata, &UserPlantsStruct)
		if err != nil {
			log.Fatal(err)
		}

		SunsynkPlantId = fmt.Sprint(UserPlantsStruct[0].Id)

		return SunsynkPlantId

	}

	return SunsynkPlantId

}

func GetInverterIDs(plantIds, token string) string {

	inverterId, err := sunsynk.GetInverterId(plantIds, token)
	if err != nil {
		log.Fatal(err)
	}

	var UserInvertersStruct sunsynk.SSApiPlantInverterDataResponse

	err = json.Unmarshal(inverterId, &UserInvertersStruct)
	if err != nil {
		log.Fatal(err)
	}

	SunsynkInverterId := UserInvertersStruct.Data.Infos[0].Sn

	return SunsynkInverterId

}
