package ixgo

import (
	"fmt"
	"go/ast"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/MeteorsLiu/llarmvp/internal/deps/pkg"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
	x "github.com/goplus/ixgo"

	xbuild "github.com/goplus/ixgo/xgobuild"
)

var ErrNoFormulaFound = fmt.Errorf("failed to get formula: no formula found")

type packageKey struct {
	PackageName string
	version.Version
}

type IXGoCompiler struct {
	ctx         *x.Context
	runnerCache map[packageKey]*Runnable
}

type formulaClassfile interface {
	Main()
}

type Runnable struct {
	Elem        any
	Dir         string
	Runner      *x.Interp
	FromVersion version.Version
	Comparator  func(a, b version.Version) int
}

func NewIXGoCompiler() *IXGoCompiler {
	return &IXGoCompiler{
		ctx:         x.NewContext(x.SupportMultipleInterp),
		runnerCache: make(map[packageKey]*Runnable),
	}
}

func (i *IXGoCompiler) FormulaOf(packageName string, packageVersion version.Version) (runner *Runnable, err error) {
	cacheKey := packageKey{packageName, packageVersion}
	runner, ok := i.runnerCache[cacheKey]
	if ok {
		return
	}

	formulaRootDir := pkg.PackagePathOf(packageName)

	var suitableRunner *Runnable
	maxCanBuild := version.None

	err = filepath.WalkDir(formulaRootDir, func(path string, d fs.DirEntry, err2 error) error {
		if err2 != nil {
			return err
		}
		if !strings.HasSuffix(path, "_llar.gox") {
			return nil
		}
		formulaDir := filepath.Dir(path)

		source, err := xbuild.BuildDir(i.ctx, formulaDir)
		if err != nil {
			return nil
		}
		pkg, err := i.ctx.LoadFile("main.go", source)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		interp, err := x.NewInterp(i.ctx, pkg)
		if err != nil {
			return nil
		}

		typ, ok := interp.GetType(structNameOf(path))
		if !ok {
			panic("cannot find struct")
		}

		val := reflect.New(typ)
		elem := val.Elem()

		val.Interface().(formulaClassfile).Main()

		fromVersion := valueOf(elem, "internalFromVersion").(version.Version)

		comparator := valueOf(elem, "onCompareFn").(func(a, b version.Version) int)

		var best bool

		if comparator(packageVersion, fromVersion) >= 0 {
			if maxCanBuild.IsNone() {
				suitableRunner = &Runnable{
					Elem:        elem,
					Dir:         formulaDir,
					Runner:      interp,
					FromVersion: fromVersion,
					Comparator:  comparator,
				}
				maxCanBuild = fromVersion
				best = true
				return nil
			}
			if comparator(fromVersion, maxCanBuild) >= 0 {
				maxCanBuild = fromVersion
				suitableRunner = &Runnable{
					Elem:        elem,
					Dir:         formulaDir,
					Runner:      interp,
					FromVersion: fromVersion,
					Comparator:  comparator,
				}
				best = true
			}
		}

		if !best {
			interp.Abort()
		}

		return nil
	})
	if suitableRunner != nil {
		runner = suitableRunner
		i.runnerCache[cacheKey] = suitableRunner
	}

	if runner == nil && err == nil {
		err = ErrNoFormulaFound
	}
	return
}

func init() {
	xbuild.RegisterClassFileType("_llar.gox", "FormulaApp", nil, "github.com/MeteorsLiu/llarmvp")
}

func structNameOf(path string) string {
	return strings.TrimSuffix(filepath.Base(path), "_llar.gox")
}

func unexportValueOf(field reflect.Value) any {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func valueOf(elem reflect.Value, name string) any {
	if ast.IsExported(name) {
		return elem.FieldByName(name).Elem().Interface()
	}
	return unexportValueOf(elem.FieldByName(name))
}
