package version

var None = PackageVersion{"none"}

type PackageVersion struct {
	Version string
}

func (v1 PackageVersion) Equal(v2 PackageVersion) bool {
	return v1.Version == v2.Version
}

func (v1 PackageVersion) IsNone() bool {
	return v1.Version == "none"
}
