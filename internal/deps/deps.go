package deps

import (
	"fmt"
	"path/filepath"

	"github.com/MeteorsLiu/llarmvp/internal/deps/pkg"
	"github.com/MeteorsLiu/llarmvp/internal/ixgo"
	"github.com/MeteorsLiu/llarmvp/internal/mvs"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

func Tidy(p *pkg.PackageDependency, ixgo *ixgo.IXGoCompiler, currentVersion string) {
	mainPackage := mvs.MvsVersion{version.From("1.7.18"), p.PackageName}

	var roots []mvs.MvsVersion

	comparatorMap := make(map[string]func(v1, v2 version.Version) int)

	mainRunner, err := ixgo.FormulaOf(p.PackageName, mainPackage.Version)
	if err != nil {
		return
	}

	comparatorMap[mainPackage.PackageName] = mainRunner.Comparator

	for _, dep := range p.Dependencies {
		// random select
		comparator, err := comparatorOf(dep, ixgo)
		if err != nil {
			return
		}
		comparatorMap[dep.PackageName] = comparator

		roots = append(roots, mvs.MvsVersion{version.From(dep.Version), dep.PackageName})
	}
	onLoad := func(mv mvs.MvsVersion) ([]mvs.MvsVersion, error) {
		subRunner, err := ixgo.FormulaOf(mv.PackageName, mv.Version)
		if err != nil {
			return nil, err
		}

		subDeps, err := pkg.Parse(filepath.Join(subRunner.Dir, "versions.json"))
		if err != nil {
			return nil, err
		}

		ret := []mvs.MvsVersion{}

		for _, dep := range subDeps.Dependencies {
			ret = append(ret, mvs.MvsVersion{version.From(dep.Version), dep.PackageName})
		}

		return ret, nil
	}
	reqs := &mvs.MvsReqs{
		Roots:         roots,
		ComparatorMap: comparatorMap,
		OnLoadVersion: onLoad,
		IsMain: func(path string, v version.Version) bool {
			return path == p.PackageName && v.Equal(mainPackage.Version)
		},
	}

	fmt.Println(mvs.Req(mainPackage, []string{mainPackage.PackageName}, reqs))
}

func comparatorOf(p pkg.Dependency, ixgo *ixgo.IXGoCompiler) (func(v1, v2 version.Version) int, error) {
	runner, err := ixgo.FormulaOf(p.PackageName, version.From(p.Version))

	if err != nil {
		return nil, err
	}

	return runner.Comparator, nil
}
