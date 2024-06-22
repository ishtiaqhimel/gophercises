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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := boltdb.SetupDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		tasks, err := boltdb.ListTasks(db)
		if err != nil {
			log.Fatal(err)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete! Why not take a vacation? 🏖")
			return
		}
		fmt.Println("You have the following tasks to do:")
		fmt.Printf(fmt.Sprintf("%%-%ds  %%20s  %%s\n", max(3, len(tasks))), "ID", "CREATED AT", "DESCRIPTION")
		format := fmt.Sprintf("%%-%dd  %%20s  %%s\n", max(3, len(tasks)))
		for i, task := range tasks {
			fmt.Printf(format, i+1, task.CreatedAt.Format("2006-01-02T15:04:05"), task.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
