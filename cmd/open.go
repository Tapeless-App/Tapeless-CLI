package cmd

import (
	"fmt"

	"tapeless.app/tapeless-cli/env"
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

			url := env.WebURL + "/projects"

			if projectInputFlag != -1 {
				url = fmt.Sprintf("%s/projects/%d", env.WebURL, projectInputFlag)
			}

			err := util.OpenBrowser(url)
			if err != nil {
				fmt.Println("Error opening browser:", err)
				return
			}
		},
	}
)
