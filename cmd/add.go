/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/chrono/task"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new task",
	Long: `Creates a new task if one with the same name does not exists already.
The new task will have the following properties:
- Name: The name of the task
- Status: The status of the task (default: pending)
- CreatedAt: The date and time the task was created
- UpdatedAt: The date and time the task was last updated
`,
	Args: cobra.ExactArgs(1),
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}

		taskRepo := task.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.CreateTask(name)
		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

