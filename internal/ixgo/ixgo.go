package ixgo

import (
	"fmt"
	"go/ast"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
	x "github.com/goplus/ixgo"
	"github.com/goplus/mod/modfile"

	xbuild "github.com/goplus/ixgo/xgobuild"
)

var ErrNoFormulaFound = fmt.Errorf("failed to get formula: no formula found")

const formulaRepo = "/Users/haolan/project/t1/llarformula"

type packageKey struct {
	PackageName string
	version.PackageVersion
}

type IXGoCompiler struct {
	ctx          *x.Context
	packageCache map[packageKey]*xbuild.Package
}

func NewIXGoCompiler() *IXGoCompiler {
	return &IXGoCompiler{ctx: x.NewContext(0), packageCache: make(map[packageKey]*xbuild.Package)}
}

func (i *IXGoCompiler) formulaOf(packageName string, packageVersion version.PackageVersion) (pkg *xbuild.Package, err error) {
	cacheKey := packageKey{packageName, packageVersion}
	pkg, ok := i.packageCache[cacheKey]
	if ok {
		return
	}
	formulaRootDir := filepath.Join(formulaRepo, packageName)

	xgoCtx := xbuild.NewContext(i.ctx)

	err = filepath.WalkDir(formulaRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, "_llar.gox") {
			return nil
		}
		currentPkg, err := xgoCtx.ParseFile(path, nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		formulaAst, err := pkg.ToAst()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		var found bool

		ast.Inspect(formulaAst, func(n ast.Node) bool {
			fmt.Println(n)
			return true
		})

		if found {
			pkg = currentPkg
			i.packageCache[cacheKey] = currentPkg
			return fs.SkipAll
		}

		return nil
	})

	if pkg == nil && err == nil {
		err = ErrNoFormulaFound
	}
	return
}

func init() {
	xbuild.RegisterProject(&modfile.Project{
		Ext:      "_llar.gox",
		Class:    "FormulaApp",
		PkgPaths: []string{"github.com/MeteorsLiu/llarmvp"},
	})
}
