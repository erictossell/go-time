package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-time/cmd"
	"go-time/db"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "go-time")
	dbFile := filepath.Join(configDir, "go-time.db")

	// Check if the directory exists, create if not
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	database, err := db.InitDB(dbFile)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer database.Close()

	var rootCmd = &cobra.Command{Use: "go-time"}
	rootCmd.AddCommand(cmd.StartCmd(database), cmd.StopCmd(database), cmd.EditCmd(database), cmd.ReadCmd(database), cmd.TuiCmd(database), cmd.DelCmd(database))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
