package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	"go-time/tui"
	// Import your TUI package
)

// TuiCmd creates a new TUI command.
func TuiCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch the Text-based User Interface",
		Run: func(cmd *cobra.Command, args []string) {
			// Call your function to start the TUI
			startTUI(db)
		},
	}
}

func startTUI(db *sql.DB) {
	// Your TUI initialization and start logic
	tui.Main(db)
}
