package build

import (
	"fmt"

	"github.com/MeteorsLiu/llarmvp"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
	"github.com/qiniu/x/stringutil"
)

type Glibc struct {
	llarmvp.FormulaApp
}

//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:4
func (this *Glibc) MainEntry() {
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:4:1
	this.FromVersion("2.0")
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:7:1
	this.PackageName__1("bminor/glibc")
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:10:1
	this.Desc__1("A massively spiffy yet delicately unobtrusive compression library")
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:13:1
	this.Homepage__1("https://github.com/bminor/glibc")
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:16:1
	this.OnBuild(func() (*llarmvp.Artifact, error) {
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:17:1
		fmt.Println(stringutil.Concat("hello libc 2.0.0 \nFormer ", this.LastArtifact().String()))
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:18:1
		return &llarmvp.Artifact{Link: func(args []string) []string {
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:19:1
			return append(args, "-Ilibc")
		}}, nil
	})
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:23:1
	this.OnSource(func(ver version.Version) (string, error) {
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:24:1
		return "", nil
	})
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:27:1
	this.OnVersions(func() []version.Version {
//line /Users/haolan/project/t1/llarformula/bminor/glibc/2.0/Glibc_llar.gox:28:1
		return nil
	})
}
func (this *Glibc) Main() {
	llarmvp.Gopt_FormulaApp_Main(this)
}
