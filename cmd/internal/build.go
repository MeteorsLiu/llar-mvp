package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build-related commands",
}

var (
	buildInfoFilter string
	buildInfoMatch  string
)

var buildInfoCmd = &cobra.Command{
	Use:   "info <PackageName>",
	Short: "Retrieve package build information",
	Long:  `Retrieve package build information with optional regex filtering.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		return fmt.Errorf("build info %s: not implemented yet (filter=%s, match=%s)",
			packageName, buildInfoFilter, buildInfoMatch)
	},
}

func init() {
	buildInfoCmd.Flags().StringVarP(&buildInfoFilter, "filter", "f", "", "Regex filter to exclude parameters")
	buildInfoCmd.Flags().StringVarP(&buildInfoMatch, "match", "m", "", "Regex filter to include parameters")
	buildCmd.AddCommand(buildInfoCmd)
	rootCmd.AddCommand(buildCmd)
}
