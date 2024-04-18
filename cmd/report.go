/*
Copyright © 2024 Mathias Donoso mathiasd88@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mathiasdonoso/chrono/db"
	"github.com/mathiasdonoso/chrono/internal/chrono/report"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Creates a report of the tasks in the database",
	Long: `Creates a report of the tasks in the database.`,
	Version: "0.0.1",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reportType := args[0]

		// TODO: Need to check how to do this in cobra
		if reportType != "daily" {
			fmt.Println("Please provide a valid report type: [daily]")
			return
		}

		dbConn, err := db.Connect()
		defer dbConn.Close()

		if err != nil {
			fmt.Println("Error connecting to db:", err)
		}

		db := dbConn.GetDB()

		reportRepository := report.NewRepository(db)
		reportService := report.NewService(reportRepository)
		reportHandler := report.NewHandler(reportService)

		res, err := reportHandler.CreateReport(reportType)

		if err != nil {
			fmt.Println("Error creating the report:", err)
			return
		}

		fmt.Println(res)
	},
}

func printReport(report report.Report) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\t Name\t Status\t Time_spent")

	// for _, t := range tasks {
	// 	id := strings.Split(t.ID, "-")[0]
	//
	// 	if len(t.Name) > 20 {
	// 		t.Name = t.Name[:20] + "..."
	// 	}
	//
	// 	fmt.Fprintln(
	// 		w,
	// 		id + "\t",
	// 		t.Name + "\t",
	// 		t.Status + "\t",
	// 		t.CreatedAt.Format("2006-01-02 15:04:05") + "\t",
	// 		t.UpdatedAt.Format("2006-01-02 15:04:05") + "\t",
	// 	)
	// }

	w.Flush()
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
