package pkg

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/MeteorsLiu/llarmvp/internal/deps"
	"github.com/MeteorsLiu/llarmvp/internal/ixgo"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

// 测试所有示例
func TestMVSExamples(t *testing.T) {
	ixgo := ixgo.NewIXGoCompiler()

	runner, err := ixgo.FormulaOf("DaveGamble/cJSON", version.Version{"1.7.18"})
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(runner.Dir)

	dep, err := deps.Parse(filepath.Join(runner.Dir, "versions.json"))
	if err != nil {
		t.Log(err)
		return
	}

	g := NewDeps(&dep)

	list, err := g.BuildList(ixgo, "1.7.18")
	fmt.Println(list, err)

	g.Require("madler/zlib", []deps.Dependency{{
		PackageName: "bminor/glibc",
		Version:     "2.42",
	}})

	list, err = g.BuildList(ixgo, "1.7.18")
	fmt.Println(list, err)

	old, ok := g.RequiredBy("DaveGamble/cJSON", version.Version{"1.7.18"})

	fmt.Println(old, ok)

	g.Require("DaveGamble/cJSON", append(old, deps.Dependency{
		PackageName: "bminor/glibc",
		Version:     "2.48",
	}))

	list, err = g.BuildList(ixgo, "1.7.18")
	fmt.Println(list, err)
}
