package time

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/prompts"
	timeService "tapeless.app/tapeless-cli/services/time"
)

var (
	removeTimeCmd = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a time entry",
		Run: func(cmd *cobra.Command, args []string) {
			project, err := GetProjectBasedOnWorkingDir("Select a project to remove time entries for:", projectIdFlag)

			if err != nil {
				fmt.Println("Aborted time entry creation", err)
				return
			}

			date, err := GetDate("For which date would you like to remove time entries?", dateFlag)

			if err != nil {
				fmt.Println("Aborted time entry creation", err)
				return
			}

			timeEntries, err := timeService.FetchTimeEntries(project.Id, date)

			if err != nil {
				fmt.Println("Failed to fetch time entries", err)
				return
			}

			if len(timeEntries) == 0 {
				fmt.Println("No time entries found for the selected date")
				return
			}

			for {
				entryToDelete, err := prompts.SelectTimeEntryPrompt("Select a time entry to remove:", timeEntries)

				if err != nil {
					fmt.Println("Aborted time entry removal", err)
					return
				}

				timeEntries, err = timeService.DeleteTimeEntry(project.Id, entryToDelete.Id)

				if err != nil {
					fmt.Println("Failed to remove time entry", err)
					return
				}

				fmt.Println("Time entry removed successfully")

				if len(timeEntries) == 0 {
					fmt.Println("No more time entries to remove")
					return
				}

				continuePrompt := promptui.Prompt{
					Label:     "Would you like to remove another time entry?",
					IsConfirm: true,
				}

				_, err = continuePrompt.Run()

				if err != nil {
					return
				}
			}

		},
	}
)

func init() {
	timeCmd.AddCommand(removeTimeCmd)
}
