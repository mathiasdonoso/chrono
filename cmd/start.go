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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start to record activity on a specific task.",
	Long: `Start recording progress on a specific task. Creates the task if not exists.`,
	Version: "0.0.1",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idOrName := args[0]

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}

		taskRepo := task.NewRepository(dbConn.GetDB())
		progressRepo := progress.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo, progressRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.StartTask(idOrName)

		// res, err := taskHandler.RemoveTaskByPartialId(id)
		if err != nil {
			fmt.Println("Error starting task:", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
