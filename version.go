package llarmvp

import "github.com/MeteorsLiu/llarmvp/pkgs/formula/version"

type VersionApp struct {
	onCompareFn func(a, b version.Version) int
}

// 提供该Package的版本比较方法，用于处理版本冲突
// 可选，当用户不提供该函数，默认使用GNU sort -V的算法
func (f *VersionApp) Compare(fn func(a, b version.Version) int) {
	f.onCompareFn = fn
}

func Gopt_VersionApp_Main(this interface{ MainEntry() }) {
	this.MainEntry()
}
