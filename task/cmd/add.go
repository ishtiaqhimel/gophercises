/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	boltdb "task/db"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := boltdb.SetupDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		des := strings.Join(args, " ")
		task, err := boltdb.AddTask(des, db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task.Description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
