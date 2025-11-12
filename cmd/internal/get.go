package internal

import (
	"github.com/MeteorsLiu/llarmvp/internal/build"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix"
	"github.com/spf13/cobra"
)

var (
	getSource bool
	getAll    bool
	getJSON   bool
)

var getCmd = &cobra.Command{
	Use:   "get <PackageName>[<PackageVersion>]",
	Short: "Retrieve package binaries or source code",
	Long:  `Retrieve package binaries or source code and output package build information.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageSpec := args[0]

		name, version, err := parsePackageSpec(packageSpec)
		if err != nil {
			return err
		}

		if version == "latest" {
			// TODO
			return nil
		}

		task, err := build.NewTask(name, matrix.Current(), version)
		if err != nil {
			return err
		}
		return task.Exec()
	},
}

func init() {
	getCmd.Flags().BoolVarP(&getSource, "source", "s", false, "Fetch only source code")
	getCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Fetch both source and binary")
	getCmd.Flags().BoolVar(&getJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(getCmd)
}
