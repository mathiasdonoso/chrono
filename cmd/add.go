/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/task"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new task",
	Long: `Creates a new task if one with the same name does not exists already.
The new task will have the following properties:
- Name: The name of the task
- Description (optional): A brief description of the task
- Status: The status of the task (default: pending)
- CreatedAt: The date and time the task was created
- UpdatedAt: The date and time the task was last updated
`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		dbConn, err := db.Connect()
		defer dbConn.Close()
		if err != nil {
			log.Fatalf("Error connecting to db: %v", err)
		}

		taskRepo := task.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.CreateTask(name, description)
		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringP("name", "n", "", "The name of the task")
	addCmd.Flags().StringP("description", "d", "", "A brief description of the task") 
}



