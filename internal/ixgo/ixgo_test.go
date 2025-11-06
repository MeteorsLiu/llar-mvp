package ixgo

import (
	"testing"

	"github.com/MeteorsLiu/llarmvp"
)

func TestFindFormula(t *testing.T) {
	ixgo := NewIXGoCompiler()

	runner, err := ixgo.formulaOf("DaveGamble/cJSON", llarmvp.Version{"1.7.18"})
	if err != nil {
		t.Log(err)
		return
	}
	runner.Elem.(interface {
		DoBuild(mrx llarmvp.Matrix) (any, error)
	}).DoBuild(llarmvp.Matrix{})
}
