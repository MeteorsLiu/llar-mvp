package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the current recipe",
	Long:  `Execute the current recipe with a root wrapper instead of direct xgo run.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
