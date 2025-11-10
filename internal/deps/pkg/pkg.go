package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/MeteorsLiu/llarmvp/internal/deps"
	"github.com/MeteorsLiu/llarmvp/internal/ixgo"
	"github.com/MeteorsLiu/llarmvp/internal/mvs"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

func newReqs(ixgo *ixgo.IXGoCompiler, p *deps.PackageDependency, currentVersion string) *mvsReqs {
	mainPackage := mvs.MvsVersion{version.From(currentVersion), p.PackageName}

	var roots []mvs.MvsVersion

	for _, dep := range p.Dependencies {
		roots = append(roots, mvs.MvsVersion{version.From(dep.Version), dep.PackageName})
	}
	onLoad := func(mv mvs.MvsVersion) (ret []mvs.MvsVersion, err error) {
		subRunner, err := ixgo.FormulaOf(mv.PackageName, mv.Version)
		if err != nil {
			return
		}
		subDeps, err := deps.Parse(filepath.Join(subRunner.Dir, "versions.json"))
		if err != nil {
			return
		}
		for _, dep := range subDeps.Dependencies {
			ret = append(ret, mvs.MvsVersion{version.From(dep.Version), dep.PackageName})
		}
		return
	}

	cmp := func(p string, v1, v2 version.Version) int {
		// fast-path: none compare, none version is always minimal.
		if v1.IsNone() && !v2.IsNone() {
			return -1
		}
		if v2.IsNone() && !v1.IsNone() {
			return 1
		}
		if v1.IsNone() && v2.IsNone() {
			return 0
		}
		return ixgo.ComparatorOf(p)(v1, v2)
	}

	return &mvsReqs{
		cmp:           cmp,
		main:          mainPackage,
		roots:         roots,
		onLoadVersion: onLoad,
		isMain: func(path string, v version.Version) bool {
			return path == p.PackageName && v.Equal(mainPackage.Version)
		},
	}
}

func Tidy(ixgo *ixgo.IXGoCompiler, p *deps.PackageDependency) {
	reqs := newReqs(ixgo, p, "")
	fmt.Println(mvs.Req(reqs.main, []string{reqs.main.PackageName}, reqs))

}

func BuildList(ixgo *ixgo.IXGoCompiler, p *deps.PackageDependency, currentVersion string) ([]deps.Dependency, error) {
	reqs := newReqs(ixgo, p, currentVersion)

	mvsDeps, err := mvs.BuildList([]mvs.MvsVersion{reqs.main}, reqs)
	if err != nil {
		return nil, err
	}

	var dependencies []deps.Dependency

	for _, dep := range mvsDeps {
		dependencies = append(dependencies, deps.Dependency{
			PackageName: dep.PackageName,
			Version:     dep.Ver,
		})
	}

	return dependencies, nil
}
