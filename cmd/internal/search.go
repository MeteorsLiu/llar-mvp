package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:   "search <Fuzzy PackageName>",
	Short: "Search packages by name",
	Long:  `Search packages by name and display package name, description, and homepage.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		return fmt.Errorf("search %s: not implemented yet (json=%v)", packageName, searchJSON)
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(searchCmd)
}
