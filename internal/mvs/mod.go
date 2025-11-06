package mvs

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"sort"
	"strings"

	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

// ModSort is like module.Sort but understands the "go" and "toolchain"
// modules and their version ordering.
func ModSort(cmp func(p string, v1, v2 version.PackageVersion) int, list []MvsVersion) {
	sort.Slice(list, func(i, j int) bool {
		mi := list[i]
		mj := list[j]
		if mi.PackageName != mj.PackageName {
			return mi.PackageName < mj.PackageName
		}
		// To help go.sum formatting, allow version/file.
		// Compare semver prefix by semver  rules,
		// file by string order.
		vi := mi.Version
		vj := mj.Version
		var fi, fj string
		if k := strings.Index(vi, "/"); k >= 0 {
			vi, fi = vi[:k], vi[k:]
		}
		if k := strings.Index(vj, "/"); k >= 0 {
			vj, fj = vj[:k], vj[k:]
		}
		if vi != vj {
			return cmp(mi.PackageName, mi.PackageVersion, mj.PackageVersion) < 0
		}
		return fi < fj
	})
}
