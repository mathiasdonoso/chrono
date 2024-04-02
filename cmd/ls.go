/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/chrono/task"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all pending tasks",
	Long: `List all task with 'pending' status and prints the result as stout
showing the following information:
- ID: The unique identifier of the task
- Name: The name of the task
- Description (optional): A brief description of the task
- Status: The status of the task (default: pending)
- CreatedAt: The date and time the task was created
- UpdatedAt: The date and time the task was last updated
`,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")

		statuses := []task.Status{task.PENDING}
		if all {
			statuses = append(statuses, task.IN_PROGRESS, task.PAUSED)
		}

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}
		taskRepo := task.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.ListTasksByStatus(statuses...)
		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}

		response := formatResponse(res)
		fmt.Println(response)
	},
}

func formatResponse(tasks []task.Task) string {
	// TODO: Make this prettier
	res := "ID\tName\tStatus\tCreated_at\tUpdated_at\n"
	for _, t := range tasks {
		i := "%s\t%s\t%s\t%s\t%s\n"

		id := strings.Split(t.ID, "-")[0]

		res += fmt.Sprintf(i,
			id,
			t.Name,
			// t.Description,
			t.Status,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			t.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return res
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("all", "a", false, "List all tasks")
}
