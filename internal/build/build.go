package build

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"

	"github.com/MeteorsLiu/llarmvp"
	"github.com/MeteorsLiu/llarmvp/internal/deps"
	"github.com/MeteorsLiu/llarmvp/internal/deps/pkg"
	"github.com/MeteorsLiu/llarmvp/internal/ixgo"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

type packageKey struct {
	PackageName    string
	PackageVersion string
}

type Task struct {
	deps           *pkg.Deps
	currentMatrix  matrix.Matrix
	currentVersion version.Version
	ixgo           *ixgo.IXGoCompiler
	formulaCache   map[packageKey]reflect.Value
}

func NewTask(packageName string, currentMatrix matrix.Matrix, currentVersion string) (*Task, error) {
	t := &Task{
		ixgo:           ixgo.NewIXGoCompiler(),
		currentVersion: version.From(currentVersion),
		currentMatrix:  currentMatrix,
		formulaCache:   make(map[packageKey]reflect.Value),
	}

	formula, err := t.ixgo.FormulaOf(packageName, t.currentVersion)
	if err != nil {
		return nil, err
	}
	dep, err := deps.Parse(filepath.Join(formula.Dir, "versions.json"))
	if err != nil {
		return nil, err
	}
	t.deps = pkg.NewDeps(&dep)

	return t, nil
}

func (t *Task) formulaOf(packageName, packageVersion string) (reflect.Value, error) {
	key := packageKey{packageName, packageVersion}
	if elem, ok := t.formulaCache[key]; ok {
		return elem, nil
	}
	formula, err := t.ixgo.FormulaOf(packageName, version.From(packageVersion))
	if err != nil {
		return reflect.Value{}, err
	}
	elem, err := formula.Elem(t.ixgo)
	if err != nil {
		return reflect.Value{}, err
	}
	t.formulaCache[key] = elem
	return elem, nil
}

func (t *Task) prepare(dir fs.FS) error {
	buildlist, err := t.deps.BuildList(t.ixgo, t.currentVersion.Ver)
	if err != nil {
		return err
	}
	main := buildlist[0]

	onRequire := func(dep deps.Dependency) error {
		formula, err := t.formulaOf(dep.PackageName, dep.Version)
		if err != nil {
			return err
		}
		// set essential value
		ixgo.SetValue(formula, "internalTempDir", dir)
		ixgo.SetValue(formula, "currentVersion", t.currentVersion)
		ixgo.SetValue(formula, "currentMatrix", t.currentMatrix)

		fn := ixgo.ValueOf(formula, "onRequireFn").(func(deps.Graph))
		if fn != nil {
			fn(t.deps)
		}
		return nil
	}

	for _, dep := range buildlist[1:] {
		if err := onRequire(dep); err != nil {
			return err
		}
	}

	return onRequire(main)
}

func (t *Task) build(dir fs.FS) error {
	buildlist, err := t.deps.BuildList(t.ixgo, t.currentVersion.Ver)
	if err != nil {
		return err
	}
	main := buildlist[0]

	var prev *llarmvp.Artifact

	onBuild := func(dep deps.Dependency) error {
		formula, err := t.formulaOf(dep.PackageName, dep.Version)
		if err != nil {
			return err
		}
		ixgo.SetValue(formula, "lastArtifact", prev)

		fn := ixgo.ValueOf(formula, "onBuildFn").(func() (result *llarmvp.Artifact, err error))

		result, err := fn()
		if err != nil {
			return err
		}
		result.BasicFormula = formula.Addr().Interface().(llarmvp.BasicFormula)

		if prev == nil {
			prev = result
		} else if result != nil {
			result.Prev = prev
			linkFn := result.Link
			result.Link = func(compileArgs []string) []string {
				return linkFn(prev.Link(compileArgs))
			}
			prev = result
		}

		return nil
	}

	for _, dep := range buildlist[1:] {
		if err := onBuild(dep); err != nil {
			return err
		}
	}

	return onBuild(main)
}

func (t *Task) Exec() error {
	dir, err := os.MkdirTemp("", "llar-build*")
	if err != nil {
		return err
	}
	fs := os.DirFS(dir)

	if err := t.prepare(fs); err != nil {
		return err
	}

	if err := t.build(fs); err != nil {
		return err
	}

	return nil
}
