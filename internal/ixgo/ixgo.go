package ixgo

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/MeteorsLiu/llarmvp"
	x "github.com/goplus/ixgo"

	xbuild "github.com/goplus/ixgo/xgobuild"
)

var ErrNoFormulaFound = fmt.Errorf("failed to get formula: no formula found")

const formulaRepo = "/Users/haolan/project/t1/llarformula"

type packageKey struct {
	PackageName string
	llarmvp.Version
}

type IXGoCompiler struct {
	ctx         *x.Context
	runnerCache map[packageKey]*Runnable
}

type buildComparable interface {
	Main()
	DoCompare(v1, v2 llarmvp.Version) int
	FromVersion__0() llarmvp.Version
}

type Runnable struct {
	Elem   any
	Runner *x.Interp
}

func NewIXGoCompiler() *IXGoCompiler {
	return &IXGoCompiler{ctx: x.NewContext(x.SupportMultipleInterp), runnerCache: make(map[packageKey]*Runnable)}
}

func (i *IXGoCompiler) formulaOf(packageName string, packageVersion llarmvp.Version) (runner *Runnable, err error) {
	cacheKey := packageKey{packageName, packageVersion}
	runner, ok := i.runnerCache[cacheKey]
	if ok {
		return
	}
	formulaRootDir := filepath.Join(formulaRepo, packageName)

	var suitableRunner *Runnable
	maxCanBuild := llarmvp.None

	err = filepath.WalkDir(formulaRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, "_llar.gox") {
			return nil
		}
		source, err := xbuild.BuildDir(i.ctx, filepath.Dir(path))
		if err != nil {
			return nil
		}
		pkg, err := i.ctx.LoadFile("main.go", source)
		if err != nil {
			return nil
		}
		interp, err := x.NewInterp(i.ctx, pkg)
		if err != nil {
			return nil
		}

		typ, ok := interp.GetType(structNameFrom(path))
		if !ok {
			panic("cannot find struct")
		}

		elem := reflect.New(typ).Interface()

		comparableElem := elem.(buildComparable)
		// init
		comparableElem.Main()

		fromVersion := comparableElem.FromVersion__0()

		var best bool

		if comparableElem.DoCompare(packageVersion, fromVersion) >= 0 {
			if maxCanBuild.IsNone() {
				suitableRunner = &Runnable{Elem: elem, Runner: interp}
				maxCanBuild = fromVersion
				best = true
				return nil
			}
			if comparableElem.DoCompare(fromVersion, maxCanBuild) >= 0 {
				maxCanBuild = fromVersion
				suitableRunner = &Runnable{Elem: elem, Runner: interp}
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

func structNameFrom(path string) string {
	return strings.TrimSuffix(filepath.Base(path), "_llar.gox")
}
