package time

import (
	"fmt"
	"strconv"

	timeService "tapeless.app/tapeless-cli/services/time"

	"github.com/spf13/cobra"
)

var (
	listTimeCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all time entries for a particular project and day",
		Run: func(cmd *cobra.Command, args []string) {

			project, err := GetProjectBasedOnWorkingDir("Select a project to create a time entry for:", projectIdFlag)

			if err != nil {
				fmt.Println(err)
				return
			}

			date, err := GetDate(fmt.Sprintf("For which date would you like to add time entries to '%s'?", project.Name), dateFlag)

			if err != nil {
				fmt.Println(err)
				return
			}

			timeEntries, err := timeService.FetchTimeEntries(project.Id, date)

			if err != nil {
				fmt.Println(err)
				return
			}

			if len(timeEntries) == 0 {
				fmt.Println("No time entries found")
				return
			}

			count := 0
			sum := 0.0

			fmt.Printf("Time entries for %s on %s:\n", project.Name, date)
			println()
			for _, timeEntry := range timeEntries {
				fmt.Printf("Time entry ID: %d\n", timeEntry.Id)
				fmt.Printf("Hours: %s\n", strconv.FormatFloat(timeEntry.Hours, 'f', -1, 64))
				fmt.Printf("Description: %s\n", timeEntry.Description)
				fmt.Println()
				count++
				sum += timeEntry.Hours
			}

			fmt.Printf("=> Total time entries: %d\n", count)
			fmt.Printf("=> Total hours: %s\n", strconv.FormatFloat(sum, 'f', -1, 64))
		},
	}
)

func init() {
	timeCmd.AddCommand(listTimeCmd)
}
