/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"
	boltdb "task/db"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete",
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

		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatal(err)
			}

			if id <= 0 && id > len(tasks) {
				log.Fatal(id, "is not a valid ID.")
			}
			ids = append(ids, id)
		}

		for _, id := range ids {
			if err := boltdb.DeleteTask(tasks[id-1].ULID, db); err != nil {
				log.Fatal(err)
			}

			ct, err := boltdb.SetTaskAsCompleted(tasks[id-1], db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Task \"%s\" is completed.\n", ct.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
