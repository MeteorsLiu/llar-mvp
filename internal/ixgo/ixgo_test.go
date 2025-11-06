package ixgo

import (
	"testing"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

func TestFindFormula(t *testing.T) {
	ixgo := NewIXGoCompiler()

	ixgo.formulaOf("DaveGamble/cJSON", version.PackageVersion{"1.7.18"})
}
