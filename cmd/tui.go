package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	"go-time/tui"
	"log"
)

// TuiCmd creates a new TUI command.
func TuiCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch the Text-based User Interface",
		Long:  "Launch the Text-based User Interface (TUI) for interactive management of timers and entries.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := startTUI(db); err != nil {
				log.Fatalf("Failed to start TUI: %v", err)
			}
		},
	}
}

func startTUI(db *sql.DB) error {
	// Your TUI initialization and start logic
	tui.Main(db)
	return nil
}
