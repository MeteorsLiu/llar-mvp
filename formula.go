package llarmvp

import (
	"io/fs"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/matrix"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
	"github.com/qiniu/x/gsh"
)

const GopPackage = true

type FormulaApp struct {
	gsh.App

	internalPackageName string
	internalDesc        string
	internalHomepage    string

	currentMatrix  matrix.PackageMatrix
	declaredMatrix matrix.PackageMatrix

	currentVersion version.PackageVersion

	onCompareFn func(a, b version.PackageVersion) int
	OnRequireFn func(fs.FS)
	onBuildFn   func(matrix matrix.PackageMatrix) (result any, err error)
	onSourceFn  func(ver version.PackageVersion) (sourceDir string, err error)
	onVersionFn func() []version.PackageVersion
}

// 返回当前PackageName
func (f *FormulaApp) PackageName__0() string {
	return f.internalPackageName
}

// 必填，声明当前LLAR Package Name，格式为：owner/repo，见下方例子
func (f *FormulaApp) PackageName__1(name string) {
	f.internalPackageName = name
}

// 返回当前描述
func (f *FormulaApp) Desc__0() string {
	return f.internalDesc
}

// 可选，添加Package Homepage页面
func (f *FormulaApp) Desc__1(desc string) {
	f.internalDesc = desc
}

// 返回当前Package Homepage URL
func (f *FormulaApp) Homepage__0() string {
	return f.internalHomepage
}

// 可选，添加Package Homepage URL
func (f *FormulaApp) Homepage__1(homepage string) {
	f.internalHomepage = homepage
}

// 返回当前Package的构建矩阵
func (f *FormulaApp) Matrix__0() matrix.PackageMatrix {
	return f.declaredMatrix
}

// 声明Package的构建矩阵
func (f *FormulaApp) Matrix__1(mrx matrix.PackageMatrix) {
	f.currentMatrix = mrx
}

// 返回当前Package的版本
func (f *FormulaApp) Version() version.PackageVersion {
	return f.currentVersion
}

// 提供该Package的版本比较方法，用于处理版本冲突
// 可选，当用户不提供该函数，默认使用GNU sort -V的算法
func (f *FormulaApp) Compare(fn func(a, b version.PackageVersion) int) {
	f.onCompareFn = fn
}

// 声明该Formula能够处理的起始版本号
func (f *FormulaApp) FromVersion(version string) {
	// do nothing
}

func (f *FormulaApp) OnRequire(fn func(dir fs.FS)) {
	f.OnRequireFn = fn
}

// 声明构建
func (f *FormulaApp) OnBuild(fn func(matrix matrix.PackageMatrix) (result any, err error)) {
	f.onBuildFn = fn
}

// 提供该Package源码下载方法，并要求维护者实现相关源码验证逻辑
func (f *FormulaApp) OnSource(fn func(ver version.PackageVersion) (sourceDir string, err error)) {
	f.onSourceFn = fn
}

// 当前配方所有版本
func (f *FormulaApp) OnVersions(fn func() []version.PackageVersion) {
	f.onVersionFn = fn
}

func Gopt_FormulaApp_Main(this interface{ MainEntry() }) {
	this.MainEntry()
}
