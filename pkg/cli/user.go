package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"ssctl/pkg/kube"
	"ssctl/pkg/sunsynk"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user data from sunksynk api",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		debugFlagValue, _ := cmd.Parent().PersistentFlags().GetBool("debug")

		if debugFlagValue {
			os.Setenv("SS_DEBUG", "TRUE")
		}

		k8sFlagValue, _ := cmd.Parent().PersistentFlags().GetBool("k8s")
		User(k8sFlagValue)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func User(k8s bool) {

	if !k8s {

		SunsynkToken := os.Getenv("SS_TOKEN")

		if SunsynkToken == "" {
			log.Fatal("No token found in env 4")
		}

		userdata, err := sunsynk.GetUserData(sunsynk.SSApiListPlantsEndpoint, SunsynkToken)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(userdata))

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

		userdata, err := sunsynk.GetUserData(sunsynk.SSApiListPlantsEndpoint, string(token))
		if err != nil {
			log.Fatal(err)
		}

		var userdatastruct sunsynk.SSApiUserPlantsResponse

		err = json.Unmarshal(userdata, &userdatastruct)
		if err != nil {
			log.Fatal(err)
		}

		result, err = kube.GetK8sSecret(clientset, "sunsynk-user-plants", "sunsynk")
		if err != nil {
			//Create the secret
			result, err = kube.CreateK8sSecret(clientset, "sunsynk-user-plants", "sunsynk", userdatastruct.Data.Infos, "plants.json")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created secret %q\n", result.GetObjectMeta().GetName())
		} else {
			result, err = kube.UpdateK8sSecret(clientset, result, "sunsynk", userdatastruct.Data.Infos, "plants.json")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Updated secret %q\n", result.GetObjectMeta().GetName())
		}

	}

}
