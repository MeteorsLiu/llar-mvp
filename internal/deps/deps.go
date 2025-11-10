package deps

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

type Dependency struct {
	PackageName string `json:"name"`
	Version     string `json:"version"`
}

type PackageDependency struct {
	PackageName  string       `json:"name"`
	Dependencies []Dependency `json:"deps"`
}

const formulaRepo = "/Users/haolan/project/t1/llarformula"

func PackagePathOf(packageName string) (dir string) {
	return filepath.Join(formulaRepo, packageName)
}

type Graph interface {
	Require(packageName string, deps []Dependency)
	RequiredBy(packageName string, version version.Version) ([]Dependency, bool)
}

func Parse(path string) (p PackageDependency, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(f).Decode(&p)
	return
}
