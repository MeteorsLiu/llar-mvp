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

	Tidy(ixgo, &dep)

	fmt.Println(BuildList(ixgo, &dep, "1.7.18"))
}
