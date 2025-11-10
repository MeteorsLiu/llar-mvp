package ixgo

import (
	"fmt"
	"go/ast"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/MeteorsLiu/llarmvp/internal/deps"
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
	ctx             *x.Context
	runnerCache     map[packageKey]*Formula
	comparatorCache map[string]func(a, b version.Version) int
}

type formulaClassfile interface {
	Main()
}

type Formula struct {
	Dir         string
	FromVersion version.Version
}

func NewIXGoCompiler() *IXGoCompiler {
	return &IXGoCompiler{
		ctx:             x.NewContext(x.SupportMultipleInterp),
		runnerCache:     make(map[packageKey]*Formula),
		comparatorCache: make(map[string]func(a version.Version, b version.Version) int),
	}
}

func (i *IXGoCompiler) comparatorOf(rootDir string) func(a, b version.Version) int {
	matches, _ := filepath.Glob(filepath.Join(rootDir, "*_version.gox"))
	if len(matches) == 0 {
		return nil
	}
	source, err := xbuild.BuildDir(i.ctx, rootDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pkg, err := i.ctx.LoadFile("main.go", source)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	interp, err := i.ctx.NewInterp(pkg)
	if err != nil {
		return nil
	}
	for _, formula := range matches {
		name := strings.TrimSuffix(filepath.Base(formula), "_version.gox")

		typ, ok := interp.GetType(name)
		if !ok {
			continue
		}

		val := reflect.New(typ)
		elem := val.Elem()

		val.Interface().(formulaClassfile).Main()

		comparator := valueOf(elem, "onCompareFn").(func(a, b version.Version) int)

		if comparator != nil {
			return comparator
		}
	}

	return nil
}

func (i *IXGoCompiler) ComparatorOf(packageName string) func(a, b version.Version) int {
	if cmp, ok := i.comparatorCache[packageName]; ok {
		return cmp
	}
	cmp := i.comparatorOf(deps.PackagePathOf(packageName))
	if cmp == nil {
		cmp = func(a, b version.Version) int {
			return version.Compare(a.Ver, b.Ver)
		}
	}
	i.comparatorCache[packageName] = cmp
	return cmp
}

func (i *IXGoCompiler) FormulaOf(packageName string, packageVersion version.Version) (runner *Formula, err error) {
	cacheKey := packageKey{packageName, packageVersion}
	runner, ok := i.runnerCache[cacheKey]
	if ok {
		return
	}

	formulaRootDir := deps.PackagePathOf(packageName)
	comparator := i.ComparatorOf(packageName)

	var suitableRunner *Formula
	maxCanBuild := version.None

	xgoCtx := xbuild.NewContext(i.ctx)

	err = filepath.WalkDir(formulaRootDir, func(path string, d fs.DirEntry, err2 error) error {
		if err2 != nil {
			return err
		}
		if !strings.HasSuffix(path, "_llar.gox") {
			return nil
		}
		formulaDir := filepath.Dir(path)

		pkg, err := xgoCtx.ParseDir(formulaDir)
		if err != nil {
			return err
		}
		formulaAst, err := pkg.ToAst()
		if err != nil {
			return err
		}

		fromVersion := version.None

		ast.Inspect(formulaAst, func(n ast.Node) bool {
			switch c := n.(type) {
			case *ast.CallExpr:
				if fn, ok := c.Fun.(*ast.SelectorExpr); ok {
					if fn.Sel.Name == "FromVersion" {
						if len(c.Args) == 0 {
							panic("invalid FromVersion argument")
						}
						arg, ok := c.Args[0].(*ast.BasicLit)
						if !ok {
							panic("invalid FromVersion argument")
						}
						fromVersion = version.From(strings.Trim(strings.Trim(arg.Value, `"`), "`"))
						return false
					}
				}
			}
			return true
		})

		if fromVersion.IsNone() {
			fmt.Printf("%s: FromVersion is not found", formulaDir)
			return nil
		}

		if comparator(packageVersion, fromVersion) >= 0 {
			if maxCanBuild.IsNone() {
				suitableRunner = &Formula{
					Dir:         formulaDir,
					FromVersion: fromVersion,
				}
				maxCanBuild = fromVersion
				return nil
			}
			if comparator(fromVersion, maxCanBuild) >= 0 {
				maxCanBuild = fromVersion
				suitableRunner = &Formula{
					Dir:         formulaDir,
					FromVersion: fromVersion,
				}
			}
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
	xbuild.RegisterClassFileType("_version.gox", "VersionApp", nil, "github.com/MeteorsLiu/llarmvp")
	xbuild.RegisterClassFileType("_llar.gox", "FormulaApp", nil, "github.com/MeteorsLiu/llarmvp")
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
