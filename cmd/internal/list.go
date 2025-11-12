package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list <PackageName>",
	Short: "List package version information",
	Long:  `List package version information including major version and specific versions.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		return fmt.Errorf("list %s: not implemented yet (json=%v)", packageName, listJSON)
	},
}

func init() {
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(listCmd)
}
