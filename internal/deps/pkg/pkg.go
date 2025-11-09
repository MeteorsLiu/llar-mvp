package pkg

import (
	"fmt"
	"path/filepath"
	"strings"

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
		// slow-path:
		// select a comparator:
		// 1. formula v1 can handle with v2, so use v1's comparator
		// 2. if case 1 dones't meet, check formula v2
		// 3. if case 1 and 2 dones't meet, that means there's a version gap between v1 and v2.
		// fallback to string compare
		runnerV1, err := ixgo.FormulaOf(p, v1)
		if err != nil {
			panic(err)
		}
		// v1 can handle
		if runnerV1.Comparator(runnerV1.FromVersion, v2) >= 0 {
			return runnerV1.Comparator(v1, v2)
		}
		runnerV2, err := ixgo.FormulaOf(p, v2)
		if err != nil {
			panic(err)
		}
		// v2 can handle
		if runnerV2.Comparator(runnerV2.FromVersion, v1) >= 0 {
			return runnerV2.Comparator(v1, v2)
		}
		// we cannot find a comparator, try to compare with two versions
		if strings.Compare(v1.Ver, v2.Ver) >= 0 {
			return runnerV1.Comparator(v1, v2)
		}
		return runnerV2.Comparator(v1, v2)
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
