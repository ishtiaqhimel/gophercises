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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove your task from the list",
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
			fmt.Printf("You have removed the \"%s\" task.\n", tasks[id-1].Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
