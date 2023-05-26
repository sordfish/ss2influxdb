package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"ssctl/pkg/sunsynk"
	"ssctl/pkg/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// inverterCmd represents the inverter command
var inverterCmd = &cobra.Command{
	Use:   "inverter",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		debugFlagValue, _ := cmd.Parent().Parent().PersistentFlags().GetBool("debug")

		if debugFlagValue {
			os.Setenv("SS_DEBUG", "TRUE")
		}

		k8sFlagValue, _ := cmd.Parent().Parent().PersistentFlags().GetBool("k8s")
		uploadFlagValue, _ := cmd.Parent().Parent().PersistentFlags().GetBool("upload")

		idata := Inverter(k8sFlagValue)

		if uploadFlagValue {
			utils.Upload2influxdb(idata)
		} else {
			fmt.Println(idata)
		}

	},
}

func init() {
	plantCmd.AddCommand(inverterCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inverterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inverterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Inverter(k8s bool) string {

	var gridRealtDataLineString, SunsynkToken, SunsynkPlantId string

	SunsynkToken = GetToken(k8s)
	SunsynkPlantId = GetPlantIDs(k8s)

	inverterId := GetInverterIDs(SunsynkPlantId, SunsynkToken)

	gridRealtimeData, err := sunsynk.GetInverterGridRealtimeData(inverterId, SunsynkToken)
	if err != nil {
		log.Fatal(err)
	}

	output, err := InverterGridRealtime2Line(SunsynkPlantId, gridRealtimeData)
	if err != nil {
		log.Fatal(err)
	}

	gridRealtDataLineString = strings.Join(output, "\n")

	return gridRealtDataLineString
}

func InverterGridRealtime2Line(plantID string, ssgridrealtimedata []byte) ([]string, error) {

	var gridrealtimedatastruct sunsynk.SSApiInverterGridRealtimeDataResponse

	var gridRealtimeDataLineStruct []LineFormat

	err := json.Unmarshal(ssgridrealtimedata, &gridrealtimedatastruct)
	if err != nil {
		return nil, err
	}

	var gridFromToday, gridToToday, gridFromTotal, gridToTotal LineFormat

	SunsynkPlantIdInt, err := strconv.Atoi(plantID)
	if err != nil {
		log.Fatal(err)
	}

	gridFromToday.PlantId = SunsynkPlantIdInt
	gridToToday.PlantId = SunsynkPlantIdInt
	gridFromTotal.PlantId = SunsynkPlantIdInt
	gridToTotal.PlantId = SunsynkPlantIdInt

	gridFromToday.Value, err = strconv.ParseFloat(gridrealtimedatastruct.Data.ETodayFrom, 32)
	if err != nil {
		log.Fatal(err)
	}

	gridToToday.Value, err = strconv.ParseFloat(gridrealtimedatastruct.Data.ETodayTo, 32)
	if err != nil {
		log.Fatal(err)
	}

	gridFromTotal.Value, err = strconv.ParseFloat(gridrealtimedatastruct.Data.ETotalFrom, 32)
	if err != nil {
		log.Fatal(err)
	}

	gridToTotal.Value, err = strconv.ParseFloat(gridrealtimedatastruct.Data.ETotalTo, 32)
	if err != nil {
		log.Fatal(err)
	}

	gridFromToday.Name = "import-today"
	gridToToday.Name = "export-today"
	gridFromTotal.Name = "import-total"
	gridToTotal.Name = "export-total"

	now := time.Now().UTC()
	epoch := now.Unix()

	gridFromToday.Timestamp = epoch
	gridToToday.Timestamp = epoch
	gridFromTotal.Timestamp = epoch
	gridToTotal.Timestamp = epoch

	gridRealtimeDataLineStruct = append(gridRealtimeDataLineStruct, gridFromToday)
	gridRealtimeDataLineStruct = append(gridRealtimeDataLineStruct, gridToToday)
	gridRealtimeDataLineStruct = append(gridRealtimeDataLineStruct, gridFromTotal)
	gridRealtimeDataLineStruct = append(gridRealtimeDataLineStruct, gridToTotal)

	var gridRealtimeDataLineStringSlice []string

	for _, line := range gridRealtimeDataLineStruct {
		gridRealtimeDataLineStringSlice = append(gridRealtimeDataLineStringSlice, "sunsynk_inverter_grid_realtime,plant="+fmt.Sprint(line.PlantId)+" "+strings.ToLower(line.Name)+"="+strconv.FormatFloat(line.Value, 'f', 2, 64)+" "+fmt.Sprint(line.Timestamp))
	}

	return gridRealtimeDataLineStringSlice, err

}
