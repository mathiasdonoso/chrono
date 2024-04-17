/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
	"github.com/mathiasdonoso/chrono/internal/chrono/task"
	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Finish a task by id.",
	Long: `Finish a task by id setting the status to "done".`,
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}

		taskRepo := task.NewRepository(dbConn.GetDB())
		progressRepo := progress.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo, progressRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.FinishTask(id)
		if err != nil {
			fmt.Println("Error finishing the task:", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
