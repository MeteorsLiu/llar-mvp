// export by github.com/goplus/ixgo/cmd/qexp

package matrix

import (
	q "github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "matrix",
		Path: "github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix",
		Deps: map[string]string{
			"runtime": "runtime",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Matrix": reflect.TypeOf((*q.Matrix)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Current": reflect.ValueOf(q.Current),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
