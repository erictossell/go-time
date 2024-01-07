package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-time/cmd"
	"go-time/db"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer database.Close()

	var rootCmd = &cobra.Command{Use: "timetracker"}
	rootCmd.AddCommand(cmd.StartCmd(database), cmd.StopCmd(database), cmd.EditCmd(database), cmd.ListCmd(database)) // Pass the database connection here
	// Add other commands as needed
	rootCmd.Execute()
}
