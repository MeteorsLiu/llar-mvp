package ixgo

import (
	"testing"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

func TestFindFormula(t *testing.T) {
	ixgo := NewIXGoCompiler()

	_, err := ixgo.FormulaOf("DaveGamble/cJSON", version.Version{"1.7.18"})
	if err != nil {
		t.Log(err)
		return
	}

}
