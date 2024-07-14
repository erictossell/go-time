package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-time/cmd"
	"go-time/db"
	"go-time/pkgs/config"
	"log"
	"os"
	"path/filepath"
)

func main() {
	config := config.New("go-time", "config.toml")
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "go-time")
	dbFile := config.Get("db_path", "go-time.db").(string)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbFilePath := filepath.Join(configDir, dbFile)
	database, err := db.InitDB(dbFilePath)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer database.Close()

	var rootCmd = &cobra.Command{
		Use:   "go-time",
		Short: "Go-Time is a time tracking application",
	}

	rootCmd.AddCommand(
		cmd.CreateCmd(database),
		cmd.StartCmd(database),
		cmd.StopCmd(database),
		cmd.EditCmd(database),
		cmd.ReadCmd(database),
		cmd.TuiCmd(database),
		cmd.DelCmd(database),
	)

	// Check if no subcommand is provided and apply command mode setting
	if len(os.Args) == 1 {
		commandMode := config.Get("command_mode", "cli").(string)
		switch commandMode {
		case "tui":
			if err := cmd.TuiCmd(database).Execute(); err != nil {
				fmt.Println("Error executing TUI command:", err)
				os.Exit(1)
			}
			return // Exit after executing TUI mode
		case "cli":
			// Default to CLI mode
		case "help":
			fmt.Println("Running in CLI mode with help page...")
			if err := rootCmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
				os.Exit(1)
			}
			return // Exit after showing help
		default:
			fmt.Println("Invalid command mode specified in config")
			os.Exit(1)
		}
	}

	// Execute the root command (handles subcommands)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}

}
