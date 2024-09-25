package cmd

import (
	"fmt"
	"io"

	"tapeless.app/tapeless-cli/util"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test command",
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := util.MakeRequest("GET", "http://localhost:4000/cli/test", nil)

		if err != nil {
			fmt.Println("Error 2:", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error 3:", err)
			return
		}

		fmt.Println("Response status:", resp.Status, string(body))
	},
}
