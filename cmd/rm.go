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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a task by id",
	Long: `Delete a task by id or any of the first 8 character of the id if there exists only one match.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}

		taskRepo := task.NewRepository(dbConn.GetDB())
		taskService := task.NewService(taskRepo)
		taskHandler := task.NewHandler(taskService)

		res, err := taskHandler.RemoveTaskByPartialId(id)
		if err != nil {
			fmt.Println("Error removing task:", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
