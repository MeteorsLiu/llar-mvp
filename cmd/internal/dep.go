package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/MeteorsLiu/llarmvp/internal/deps"
	"github.com/spf13/cobra"
)

var _depRegexp = regexp.MustCompile(`(.*/.*)(\[.*\])?`)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Dependency management commands",
}

var depInitCmd = &cobra.Command{
	Use:   "init <PackageName>",
	Short: "Initialize dependency module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		if _, err := os.Stat("versions.json"); err == nil {
			return nil
		}

		p := deps.PackageDependency{
			PackageName: packageName,
		}

		f, err := os.Create("versions.json")
		if err != nil {
			return err
		}

		return json.NewEncoder(f).Encode(&p)
	},
}

var depGetCmd = &cobra.Command{
	Use:   "get <PackageName>[<PackageVersion>]",
	Short: "Add dependency to current package",
	Long:  `Add dependency to current package. Example: llar dep get madler/zlib[1.2.1]`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageSpec := args[0]

		if _, err := os.Stat("versions.json"); os.IsNotExist(err) {
			return fmt.Errorf("dep get %s: mod is not init", packageSpec)
		}

		dep, err := deps.Parse("versions.json")
		if err != nil {
			return err
		}

		name, version, err := parsePackageSpec(packageSpec)
		if err != nil {
			return err
		}
		if version == "latest" {
			// TODO(MeteorsLiu): look up latest version
			return nil
		}

		dep.Dependencies = append(dep.Dependencies, deps.Dependency{
			PackageName: name,
			Version:     version,
		})

		f, err := os.OpenFile("versions.json", os.O_RDWR, 0700)
		if err != nil {
			return err
		}

		return json.NewEncoder(f).Encode(&dep)
	},
}

var depTidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Organize current dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("dep tidy: not implemented yet")
	},
}

var depListCmd = &cobra.Command{
	Use:   "list <PackageName>",
	Short: "List package dependencies",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		return fmt.Errorf("dep list %s: not implemented yet", packageName)
	},
}

func parsePackageSpec(spec string) (packageName, packageVersion string, err error) {
	ret := _depRegexp.FindAllStringSubmatch(spec, -1)

	switch len(ret) {
	case 1:
		packageName = ret[0][0]
		packageVersion = "latest"
	case 2:
		packageName = ret[0][0]
		packageVersion = ret[1][0]
	default:
		err = fmt.Errorf("failed to parse input: %s", spec)
	}

	return
}

func init() {
	depCmd.AddCommand(depInitCmd)
	depCmd.AddCommand(depGetCmd)
	depCmd.AddCommand(depTidyCmd)
	depCmd.AddCommand(depListCmd)
	rootCmd.AddCommand(depCmd)
}
