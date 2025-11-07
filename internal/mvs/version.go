package mvs

import (
	"github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

type MvsVersion struct {
	version.Version
	PackageName string
}

// MvsReqs implements mvs.Reqs for module semantic versions,
// with any exclusions or replacements applied internally.
type MvsReqs struct {
	Roots         []MvsVersion
	IsMain        func(path string, v version.Version) bool
	ComparatorMap map[string]func(v1, v2 version.Version) int
	OnLoadVersion func(MvsVersion) ([]MvsVersion, error)
}

func (r *MvsReqs) Required(mod MvsVersion) ([]MvsVersion, error) {
	if r.IsMain(mod.PackageName, mod.Version) {
		// Use the build list as it existed when r was constructed, not the current
		// global build list.
		return r.Roots, nil
	}

	if mod.Version.IsNone() {
		return nil, nil
	}

	return r.OnLoadVersion(mod)
}

// Max returns the maximum of v1 and v2 according to gover.ModCompare.
//
// As a special case, the version "" is considered higher than all other
// versions. The main module (also known as the target) has no version and must
// be chosen over other versions of the same module in the module dependency
// graph.
func (r *MvsReqs) Max(p string, v1, v2 version.Version) version.Version {
	if r.cmpVersion(p, v1, v2) < 0 {
		return v2
	}
	return v1
}

// Upgrade is a no-op, here to implement mvs.Reqs.
// The upgrade logic for go get -u is in ../modget/get.go.
func (*MvsReqs) Upgrade(m MvsVersion) (MvsVersion, error) {
	return m, nil
}

// cmpVersion implements the comparison for versions in the module loader.
//
// It is consistent with gover.ModCompare except that as a special case,
// the version "" is considered higher than all other versions.
// The main module (also known as the target) has no version and must be chosen
// over other versions of the same module in the module dependency graph.
func (m *MvsReqs) cmpVersion(p string, v1, v2 version.Version) int {
	if m.IsMain(p, v2) {
		if m.IsMain(p, v1) {
			return 0
		}
		return -1
	}
	if m.IsMain(p, v1) {
		return 1
	}
	return m.ComparatorMap[p](v1, v2)
}
