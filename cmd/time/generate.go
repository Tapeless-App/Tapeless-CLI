package time

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	sync "tapeless.app/tapeless-cli/cmd/sync"
	syncService "tapeless.app/tapeless-cli/services/sync"
	timeService "tapeless.app/tapeless-cli/services/time"
)

var (
	generateTimeCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "ai"},
		Short:   "Generate time entries based on synced commits for a project using AI",
		Run: func(cmd *cobra.Command, args []string) {
			project, err := GetProjectBasedOnWorkingDir("Select a project to generate time entries for:", projectIdFlag)

			if err != nil {
				fmt.Println(err)
				return
			}

			date, err := GetDate(fmt.Sprintf("Which date would you like to generate time entries for '%s'?", project.Name), dateFlag)

			if err != nil {
				fmt.Println(err)
				return
			}

			timeEntries, err := timeService.FetchTimeEntries(project.Id, date)

			if err != nil {
				fmt.Println(err)
				return
			}

			count := 0
			sum := 0.0

			for _, timeEntry := range timeEntries {
				count++
				sum += timeEntry.Hours
			}

			// Ensure that commits are synced before generating time entries
			for {
				commits, err := syncService.FetchCommitsForProjectAndDate(project.Id, date)

				if err != nil {
					fmt.Println("Error fetching commits", err)
					return
				}

				label := ""
				items := []string{"Sync"}

				if len(commits) == 0 {
					label = "No commits found - please sync your commits first before continuing"
				} else {
					label = fmt.Sprintf("Found %d commits for %s on %s. Would you like to sync before moving on?", len(commits), project.Name, date)
					items = append(items, "Generate")
				}

				items = append(items, "Abort")

				nextStepPrompt := promptui.Select{
					Label: label,
					Items: items,
					Size:  len(items),
				}

				_, result, err := nextStepPrompt.Run()

				if err != nil {
					fmt.Println(err)
					return
				}

				if result == "Abort" {
					fmt.Println("Time generation aborted")
					return
				}

				if result == "Generate" {
					break
				}

				if result == "Sync" {
					fmt.Println("Syncing commits, this might take a few seconds...")

					sync.SyncCmd.Run(cmd, args)

					fmt.Println("Successfully synced commits")
				}
			}

			label := ""

			if sum > 0 {
				label = fmt.Sprintf("You currently have %s hours logged for %s on %s. ", strconv.FormatFloat(sum, 'f', -1, 64), project.Name, date)
			}

			targetHoursPrompt := promptui.Prompt{
				Label: label + fmt.Sprintf("How many hours did you work in total for %s on %s?", project.Name, date),
				Validate: func(input string) error {
					num, err := strconv.ParseFloat(input, 32)
					if err != nil {
						return fmt.Errorf("could not convert %s to float", input)
					}
					if num <= sum {
						return fmt.Errorf("to generate time entries, your total hours must be greater than your existing time entries")
					}
					return nil
				},
			}

			targetHoursStr, err := targetHoursPrompt.Run()

			if err != nil {
				fmt.Println(err)
				return
			}

			targetHours, err := strconv.ParseFloat(targetHoursStr, 64)

			if err != nil {
				fmt.Println(err)
				return
			}

			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("AI will generated time entries to fill up the remaining %s hours based on your commits. Continue", strconv.FormatFloat(targetHours-sum, 'f', -1, 64)),
				IsConfirm: true,
				Default:   "y",
			}

			_, err = confirmPrompt.Run()

			if err != nil {
				fmt.Println("Time generation aborted")
				return
			}

			fmt.Println("Generating time entries, this might take a few seconds...")

			timeEntries, err = timeService.GenerateTimeEntries(project.Id, date, targetHours-sum)

			if err != nil {
				fmt.Println(err)
				return
			}

			if len(timeEntries) == 0 {
				fmt.Println("No time entries were generated")
				return
			}

			fmt.Printf("Successfully generated time entries for %s on %s.\n", project.Name, date)

			count = 0
			sum = 0.0

			fmt.Printf("Your new time entries for %s on %s:\n", project.Name, date)
			println()
			for _, timeEntry := range timeEntries {
				fmt.Printf("Time entry ID: %d\n", timeEntry.Id)
				fmt.Printf("Hours: %s\n", strconv.FormatFloat(timeEntry.Hours, 'f', -1, 64))
				fmt.Printf("Description: %s\n", timeEntry.Description)
				fmt.Println()
				count++
				sum += timeEntry.Hours
			}

			fmt.Printf("=> New total time entries: %d\n", count)
			fmt.Printf("=> New total hours: %s\n", strconv.FormatFloat(sum, 'f', -1, 64))

		}}
)

func init() {
	timeCmd.AddCommand(generateTimeCmd)
}
