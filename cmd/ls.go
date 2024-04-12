/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
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
	Version: "0.0.1",
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
		progressRepo := progress.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo, progressRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.ListTasksByStatus(statuses...)
		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}

		printFormatedResponse(res)
	},
}

func printFormatedResponse(tasks []task.Task) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\t Name\t Status\t Created_at\t Updated_at")

	for _, t := range tasks {
		id := strings.Split(t.ID, "-")[0]

		if len(t.Name) > 20 {
			t.Name = t.Name[:20] + "..."
		}

		fmt.Fprintln(
			w,
			id + "\t",
			t.Name + "\t",
			t.Status + "\t",
			t.CreatedAt.Format("2006-01-02 15:04:05") + "\t",
			t.UpdatedAt.Format("2006-01-02 15:04:05") + "\t",
		)
	}

	w.Flush()
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("all", "a", false, "List all tasks")
}
