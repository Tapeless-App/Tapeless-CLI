package time

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	timeService "tapeless.app/tapeless-cli/services/time"
)

var (
	addTimeCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"create", "new"},
		Short:   "Add time entry for a project",
		Long:    "Add time entry for a project. If working directory is in a repo that is added to a project, will default to that project.",
		Run: func(cmd *cobra.Command, args []string) {

			project, err := GetProjectBasedOnWorkingDir("Select a project to create a time entry for:", projectIdFlag)

			if err != nil {
				fmt.Println("Aborted time entry creation", err)
				return
			}

			date, err := GetDate(fmt.Sprintf("For which date would you like to add time entries for '%s'?", project.Name), dateFlag)

			if err != nil {
				fmt.Println("Aborted time entry creation", err)
				return
			}

			units := "hours" // future feature - make this configurable

			for {

				timePrompt := promptui.Prompt{
					Label: fmt.Sprintf("How many %s do you wish to add?", units),
					Validate: func(input string) error {
						num, err := strconv.ParseFloat(input, 32)
						if err != nil {
							return fmt.Errorf("could not convert %s to float", input)
						}
						if num <= 0 {
							return fmt.Errorf("time entry must be greater than 0")
						}
						return nil
					},
				}

				timeStr, err := timePrompt.Run()

				if err != nil {
					fmt.Println("Time creation aborted", err)
					return
				}

				time, err := strconv.ParseFloat(timeStr, 64)

				if err != nil {
					fmt.Println(err)
					return
				}

				descriptionPrompt := promptui.Prompt{
					Label: "(Optional) Description - what did you work on?",
				}

				description, err := descriptionPrompt.Run()

				if err != nil {
					fmt.Println("Time creation aborted", err)
				}

				timeEntries, err := timeService.CreateTimeEntry(project.Id, timeService.TimeEntryCreateRequest{
					Description: description,
					Date:        date,
					Hours:       time,
				})

				if err != nil {
					fmt.Println("Error creating time entry - aborting", err)
					return
				}

				fmt.Printf("Successfully added %s %s to project '%s'.\n",
					strconv.FormatFloat(time, 'f', -1, 64), units, project.Name)

				fmt.Printf("Total entries for '%s' on %s: %d entries, %s hrs.\n",
					project.Name, date, timeEntries.TimeEntriesCount, strconv.FormatFloat(timeEntries.TotalHours, 'f', -1, 64))

				addEntryPrompt := promptui.Prompt{
					Label:     fmt.Sprintf("Would you like to add another time entry to %s on %s", project.Name, date),
					Default:   "y",
					IsConfirm: true,
					AllowEdit: true,
				}

				_, err = addEntryPrompt.Run()

				if err != nil {
					fmt.Println("Adding time entries completed")
					return
				}
			}

		},
	}
)

func init() {

	timeCmd.AddCommand(addTimeCmd)

}
