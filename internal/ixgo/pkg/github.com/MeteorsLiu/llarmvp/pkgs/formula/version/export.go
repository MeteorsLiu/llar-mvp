// export by github.com/goplus/ixgo/cmd/qexp

package version

import (
	q "github.com/MeteorsLiu/llarmvp/pkgs/formula/version"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name:       "version",
		Path:       "github.com/MeteorsLiu/llarmvp/pkgs/formula/version",
		Deps:       map[string]string{},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Version": reflect.TypeOf((*q.Version)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"None": reflect.ValueOf(&q.None),
		},
		Funcs: map[string]reflect.Value{
			"Compare": reflect.ValueOf(q.Compare),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
