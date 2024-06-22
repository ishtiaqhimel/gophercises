/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	boltdb "task/db"

	"github.com/spf13/cobra"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := boltdb.SetupDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		tasks, err := boltdb.ListCompletedTasks(db)
		if err != nil {
			log.Fatal(err)
		}
		if len(tasks) == 0 {
			fmt.Println("You have not completed any task yet!")
			return
		}
		fmt.Println("You have completed the following tasks so far:")
		fmt.Printf(fmt.Sprintf("%%-%ds  %%20s  %%s\n", max(3, len(tasks))), "ID", "COMPLETED AT", "DESCRIPTION")
		format := fmt.Sprintf("%%-%dd  %%20s  %%s\n", max(3, len(tasks)))
		for i, task := range tasks {
			fmt.Printf(format, i+1, task.CompletedAt.Format("2006-01-02T15:04:05"), task.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
