package pkg

import (
	"path/filepath"

	"github.com/MeteorsLiu/llarmvp/internal/deps"
	"github.com/MeteorsLiu/llarmvp/internal/ixgo"
	"github.com/MeteorsLiu/llarmvp/internal/mvs"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

type Deps struct {
	Pkg         *deps.PackageDependency
	Graph       *mvs.Graph
	cacheResult []deps.Dependency
	override    map[string][]deps.Dependency
}

func NewDeps(p *deps.PackageDependency) *Deps {
	return &Deps{Pkg: p, override: make(map[string][]deps.Dependency)}
}

func (p *Deps) Require(packageName string, deps []deps.Dependency) {
	p.override[packageName] = deps
	p.cacheResult = nil
}

func (p *Deps) RequiredBy(packageName string, version version.Version) ([]deps.Dependency, bool) {
	reqs, ok := p.Graph.RequiredBy(mvs.MvsVersion{version, packageName})
	if !ok {
		return nil, ok
	}
	var ret []deps.Dependency

	for _, dep := range reqs {
		ret = append(ret, deps.Dependency{PackageName: dep.PackageName, Version: dep.Ver})
	}

	return ret, true
}

func formulaPackageNameOf(ixgo *ixgo.IXGoCompiler, packageName string, packageVersion version.Version) string {
	runner, err := ixgo.FormulaOf(packageName, packageVersion)
	if err != nil {
		panic(err)
	}
	return runner.PackageName
}

func (p *Deps) reqs(ixgo *ixgo.IXGoCompiler, currentVersion string) *mvsReqs {
	ver := version.From(currentVersion)
	mainPackage := mvs.MvsVersion{ver, formulaPackageNameOf(ixgo, p.Pkg.PackageName, ver)}

	var roots []mvs.MvsVersion

	mainDeps := p.override[p.Pkg.PackageName]

	if mainDeps == nil {
		mainDeps = p.Pkg.Dependencies
	}

	for _, dep := range mainDeps {
		depVer := version.From(dep.Version)
		roots = append(roots, mvs.MvsVersion{depVer, formulaPackageNameOf(ixgo, dep.PackageName, depVer)})
	}

	onLoad := func(mv mvs.MvsVersion) (ret []mvs.MvsVersion, err error) {
		if deps, ok := p.override[mv.PackageName]; ok {
			for _, dep := range deps {
				depVer := version.From(dep.Version)
				ret = append(ret, mvs.MvsVersion{depVer, formulaPackageNameOf(ixgo, dep.PackageName, depVer)})
			}
			return
		}
		formula, err := ixgo.FormulaOf(mv.PackageName, mv.Version)
		if err != nil {
			return
		}
		subDeps, err := deps.Parse(filepath.Join(formula.Dir, "versions.json"))
		if err != nil {
			return
		}
		for _, dep := range subDeps.Dependencies {
			depVer := version.From(dep.Version)
			ret = append(ret, mvs.MvsVersion{depVer, formulaPackageNameOf(ixgo, dep.PackageName, depVer)})
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
			return path == mainPackage.PackageName && v.Equal(mainPackage.Version)
		},
	}
}

func Tidy(ixgo *ixgo.IXGoCompiler, p *deps.PackageDependency, currentVersion string) {
	// reqs := newReqs(ixgo, p, currentVersion)

}

func (d *Deps) BuildList(ixgo *ixgo.IXGoCompiler, currentVersion string) ([]deps.Dependency, error) {
	if d.cacheResult != nil {
		return d.cacheResult, nil
	}
	var err error
	d.Graph, d.cacheResult, err = d.buildList(ixgo, currentVersion)

	return d.cacheResult, err
}

func (d *Deps) buildList(ixgo *ixgo.IXGoCompiler, currentVersion string) (*mvs.Graph, []deps.Dependency, error) {
	reqs := d.reqs(ixgo, currentVersion)

	graph, mvsDeps, err := mvs.BuildList([]mvs.MvsVersion{reqs.main}, reqs)
	if err != nil {
		return nil, nil, err
	}

	var dependencies []deps.Dependency

	for _, dep := range mvsDeps {
		dependencies = append(dependencies, deps.Dependency{
			PackageName: dep.PackageName,
			Version:     dep.Ver,
		})
	}

	return graph, dependencies, nil
}
