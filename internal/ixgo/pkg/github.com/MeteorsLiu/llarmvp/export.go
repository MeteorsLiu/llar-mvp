// export by github.com/goplus/ixgo/cmd/qexp

package llarmvp

import (
	q "github.com/MeteorsLiu/llarmvp"

	"go/constant"
	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "llarmvp",
		Path: "github.com/MeteorsLiu/llarmvp",
		Deps: map[string]string{
			"github.com/MeteorsLiu/llarmvp/internal/deps":        "deps",
			"github.com/MeteorsLiu/llarmvp/pkgs/formula/gsh":     "gsh",
			"github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix":  "matrix",
			"github.com/MeteorsLiu/llarmvp/pkgs/formula/version": "version",
			"io/fs": "fs",
		},
		Interfaces: map[string]reflect.Type{
			"BasicFormula": reflect.TypeOf((*q.BasicFormula)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"Artifact":   reflect.TypeOf((*q.Artifact)(nil)).Elem(),
			"FormulaApp": reflect.TypeOf((*q.FormulaApp)(nil)).Elem(),
			"VersionApp": reflect.TypeOf((*q.VersionApp)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Gopt_FormulaApp_Main": reflect.ValueOf(q.Gopt_FormulaApp_Main),
			"Gopt_VersionApp_Main": reflect.ValueOf(q.Gopt_VersionApp_Main),
		},
		TypedConsts: map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{
			"GopPackage": {"untyped bool", constant.MakeBool(bool(q.GopPackage))},
		},
	})
}
