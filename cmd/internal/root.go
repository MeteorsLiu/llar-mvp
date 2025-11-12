package internal

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "llar",
	Short: "LLAR - A C/C++ package manager",
	Long:  `LLAR is a package manager for C/C++ projects with centralized repository management.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
