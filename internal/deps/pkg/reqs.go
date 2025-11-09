package pkg

import (
	"github.com/MeteorsLiu/llarmvp/internal/mvs"
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

// mvsReqs implements mvs.Reqs for module semantic versions,
// with any exclusions or replacements applied internally.
type mvsReqs struct {
	main          mvs.MvsVersion
	roots         []mvs.MvsVersion
	isMain        func(path string, v version.Version) bool
	cmp           func(p string, v1, v2 version.Version) int
	onLoadVersion func(mvs.MvsVersion) ([]mvs.MvsVersion, error)
}

func (r *mvsReqs) Required(mod mvs.MvsVersion) ([]mvs.MvsVersion, error) {
	if r.isMain(mod.PackageName, mod.Version) {
		// Use the build list as it existed when r was constructed, not the current
		// global build list.
		return r.roots, nil
	}

	if mod.Version.IsNone() {
		return nil, nil
	}

	return r.onLoadVersion(mod)
}

// Max returns the maximum of v1 and v2 according to gover.ModCompare.
//
// As a special case, the version "" is considered higher than all other
// versions. The main module (also known as the target) has no version and must
// be chosen over other versions of the same module in the module dependency
// graph.
func (r *mvsReqs) Max(p string, v1, v2 version.Version) version.Version {
	if r.cmpVersion(p, v1, v2) < 0 {
		return v2
	}
	return v1
}

// Upgrade is a no-op, here to implement mvs.Reqs.
// The upgrade logic for go get -u is in ../modget/get.go.
func (*mvsReqs) Upgrade(m mvs.MvsVersion) (mvs.MvsVersion, error) {
	return m, nil
}

// cmpVersion implements the comparison for versions in the module loader.
//
// It is consistent with gover.ModCompare except that as a special case,
// the version "" is considered higher than all other versions.
// The main module (also known as the target) has no version and must be chosen
// over other versions of the same module in the module dependency graph.
func (m *mvsReqs) cmpVersion(p string, v1, v2 version.Version) int {
	if m.isMain(p, v2) {
		if m.isMain(p, v1) {
			return 0
		}
		return -1
	}
	if m.isMain(p, v1) {
		return 1
	}
	return m.cmp(p, v1, v2)
}
