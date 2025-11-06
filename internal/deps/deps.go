package deps

type Dependency struct {
	PackageName string `json:"name"`
	Version     string `json:"version"`
}

type PackageDependency struct {
	PackageName  string       `json:"name"`
	Dependencies []Dependency `json:"deps"`
}

func (p *PackageDependency) buildGraph(currentVersion string) {

}
