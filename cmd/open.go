package cmd

import (
	"fmt"

	"tapeless.app/tapeless-cli/util"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(openCmd)
	openCmd.Flags().IntVarP(&projectInputFlag, "project-id", "p", -1, "Project ID")
}

var (
	projectInputFlag int
	openCmd          = &cobra.Command{
		Use:   "open",
		Short: "Open the Tapeless web app",
		Run: func(cmd *cobra.Command, args []string) {

			url := "http://localhost:5173/projects"

			if projectInputFlag != -1 {
				url = fmt.Sprintf("http://localhost:5173/projects/%d", projectInputFlag)
			}

			err := util.OpenBrowser(url)
			if err != nil {
				fmt.Println("Error opening browser:", err)
				return
			}
		},
	}
)
