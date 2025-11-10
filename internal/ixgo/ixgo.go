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
	pkg         *xbuild.Package
	structName  string
	Dir         string
	PackageName string
	FromVersion version.Version
}

func NewIXGoCompiler() *IXGoCompiler {
	i := &IXGoCompiler{
		ctx:             x.NewContext(x.SupportMultipleInterp),
		runnerCache:     make(map[packageKey]*Formula),
		comparatorCache: make(map[string]func(a version.Version, b version.Version) int),
	}
	return i
}

func (i *IXGoCompiler) comparatorOf(rootDir string) func(a, b version.Version) int {
	matches, _ := filepath.Glob(filepath.Join(rootDir, "*_version.gox"))
	if len(matches) == 0 {
		return nil
	}
	lookup := i.ctx.Lookup
	defer func() {
		i.ctx.Lookup = lookup
	}()
	i.ctx.Lookup = func(_, path string) (dir string, found bool) {
		return newGoModDriver().Lookup(rootDir, path)
	}
	source, err := xbuild.BuildDir(i.ctx, rootDir)
	if err != nil {
		return nil
	}
	pkgs, err := i.ctx.LoadFile("main.go", source)
	if err != nil {
		return nil
	}

	interp, err := i.ctx.NewInterp(pkgs)
	if err != nil {
		return nil
	}
	interp.RunInit()

	for _, formula := range matches {
		name := strings.TrimSuffix(filepath.Base(formula), "_version.gox")

		typ, ok := interp.GetType(name)
		if !ok {
			continue
		}

		val := reflect.New(typ)
		elem := val.Elem()

		val.Interface().(formulaClassfile).Main()

		comparator := ValueOf(elem, "onCompareFn").(func(a, b version.Version) int)

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
		fmt.Printf("%s use default comparator\n", packageName)
		cmp = func(a, b version.Version) int {
			return version.Compare(a.Ver, b.Ver)
		}
	}
	i.comparatorCache[packageName] = cmp
	return cmp
}

func (f *Formula) Elem(ixgo *IXGoCompiler) (reflect.Value, error) {
	source, err := f.pkg.ToSource()
	if err != nil {
		return reflect.Value{}, err
	}
	pkg, err := ixgo.ctx.LoadFile("main.go", source)
	if err != nil {
		return reflect.Value{}, err

	}
	interp, err := ixgo.ctx.NewInterp(pkg)
	if err != nil {
		return reflect.Value{}, err

	}
	err = interp.RunInit()
	if err != nil {
		return reflect.Value{}, err
	}
	typ, ok := interp.GetType(f.structName)
	if !ok {
		return reflect.Value{}, fmt.Errorf("failed to get struct name: %s", f.structName)
	}
	elem := reflect.New(typ)
	elem.Interface().(formulaClassfile).Main()

	return elem.Elem(), nil
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

		var structName string
		var packageName string
		fromVersion := version.None

		ast.Inspect(formulaAst, func(n ast.Node) bool {
			switch c := n.(type) {
			case *ast.CallExpr:
				switch fn := c.Fun.(type) {
				case *ast.SelectorExpr:
					switch fn.Sel.Name {
					case "FromVersion":
						var ver string
						ver, err = parseCallArg(packageName, fn.Sel.Name, c)
						if err != nil {
							return false
						}
						fromVersion = version.From(ver)
					case "PackageName__1":
						var formulaPkgName string
						formulaPkgName, err = parseCallArg(packageName, fn.Sel.Name, c)
						if err != nil {
							return false
						}
						packageName = formulaPkgName
					}
				case *ast.Ident:
					if fn.Name == "new" {
						structName, err = parseCallArg(packageName, fn.Name, c)
						if err != nil {
							return false
						}
					}
				}
			}
			return true
		})

		if err != nil {
			return fs.SkipAll
		}

		if fromVersion.IsNone() {
			fmt.Printf("%s: FromVersion is not found", formulaDir)
			return nil
		}

		if comparator(packageVersion, fromVersion) >= 0 {
			if maxCanBuild.IsNone() {
				suitableRunner = &Formula{
					pkg:         pkg,
					structName:  structName,
					Dir:         formulaDir,
					PackageName: packageName,
					FromVersion: fromVersion,
				}
				maxCanBuild = fromVersion
				return nil
			}
			if comparator(fromVersion, maxCanBuild) >= 0 {
				maxCanBuild = fromVersion
				suitableRunner = &Formula{
					pkg:         pkg,
					structName:  structName,
					Dir:         formulaDir,
					PackageName: packageName,
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

func unexportValueOf(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

func ValueOf(elem reflect.Value, name string) any {
	if ast.IsExported(name) {
		return elem.FieldByName(name).Elem().Interface()
	}
	return unexportValueOf(elem.FieldByName(name)).Interface()
}

func SetValue(elem reflect.Value, name string, value any) {
	if ast.IsExported(name) {
		elem.FieldByName(name).Elem().Set(reflect.ValueOf(value))
	}
	unexportValueOf(elem.FieldByName(name)).Set(reflect.ValueOf(value))
}

func parseCallArg(pkgName, fnName string, c *ast.CallExpr) (string, error) {
	if len(c.Args) == 0 {
		return "", fmt.Errorf("%s invalid argument: %s", pkgName, fnName)
	}
	var argResult string
	switch arg := c.Args[0].(type) {
	case *ast.BasicLit:
		argResult = strings.Trim(strings.Trim(arg.Value, `"`), "`")
		if argResult == "" {
			return "", fmt.Errorf("%s empty args: %s", pkgName, fnName)
		}
	case *ast.Ident:
		argResult = arg.Name
		if argResult == "" {
			return "", fmt.Errorf("%s empty args: %s", pkgName, fnName)
		}
	}

	return argResult, nil
}
