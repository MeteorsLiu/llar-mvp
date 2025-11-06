package mvs

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"slices"

	pv "github.com/MeteorsLiu/llarmvp/pkgs/formula/version"
)

// Graph implements an incremental version of the MVS algorithm, with the
// requirements pushed by the caller instead of pulled by the MVS traversal.
type Graph struct {
	cmp   func(p string, v1, v2 pv.PackageVersion) int
	roots []MvsVersion

	required map[MvsVersion][]MvsVersion

	isRoot   map[MvsVersion]bool          // contains true for roots and false for reachable non-roots
	selected map[string]pv.PackageVersion // path → version
}

// NewGraph returns an incremental MVS graph containing only a set of root
// dependencies and using the given max function for version strings.
//
// The caller must ensure that the root slice is not modified while the Graph
// may be in use.
func NewGraph(cmp func(p string, v1, v2 pv.PackageVersion) int, roots []MvsVersion) *Graph {
	g := &Graph{
		cmp:      cmp,
		roots:    slices.Clip(roots),
		required: make(map[MvsVersion][]MvsVersion),
		isRoot:   make(map[MvsVersion]bool),
		selected: make(map[string]pv.PackageVersion),
	}

	for _, m := range roots {
		g.isRoot[m] = true
		if g.cmp(m.PackageName, g.Selected(m.PackageName), m.PackageVersion) < 0 {
			g.selected[m.PackageName] = m.PackageVersion
		}
	}

	return g
}

// Require adds the information that module m requires all modules in reqs.
// The reqs slice must not be modified after it is passed to Require.
//
// m must be reachable by some existing chain of requirements from g's target,
// and Require must not have been called for it already.
//
// If any of the modules in reqs has the same path as g's target,
// the target must have higher precedence than the version in req.
func (g *Graph) Require(m MvsVersion, reqs []MvsVersion) {
	// To help catch disconnected-graph bugs, enforce that all required versions
	// are actually reachable from the roots (and therefore should affect the
	// selected versions of the modules they name).
	if _, reachable := g.isRoot[m]; !reachable {
		panic(fmt.Sprintf("%v is not reachable from any root", m))
	}

	// Truncate reqs to its capacity to avoid aliasing bugs if it is later
	// returned from RequiredBy and appended to.
	reqs = slices.Clip(reqs)

	if _, dup := g.required[m]; dup {
		panic(fmt.Sprintf("requirements of %v have already been set", m))
	}
	g.required[m] = reqs

	for _, dep := range reqs {
		// Mark dep reachable, regardless of whether it is selected.
		if _, ok := g.isRoot[dep]; !ok {
			g.isRoot[dep] = false
		}

		if g.cmp(dep.PackageName, g.Selected(dep.PackageName), dep.PackageVersion) < 0 {
			g.selected[dep.PackageName] = dep.PackageVersion
		}
	}
}

// RequiredBy returns the slice of requirements passed to Require for m, if any,
// with its capacity reduced to its length.
// If Require has not been called for m, RequiredBy(m) returns ok=false.
//
// The caller must not modify the returned slice, but may safely append to it
// and may rely on it not to be modified.
func (g *Graph) RequiredBy(m MvsVersion) (reqs []MvsVersion, ok bool) {
	reqs, ok = g.required[m]
	return reqs, ok
}

// Selected returns the selected version of the given module path.
//
// If no version is selected, Selected returns version "none".
func (g *Graph) Selected(path string) (version pv.PackageVersion) {
	v, ok := g.selected[path]
	if !ok {
		return pv.None
	}
	return v
}

// BuildList returns the selected versions of all modules present in the Graph,
// beginning with the selected versions of each module path in the roots of g.
//
// The order of the remaining elements in the list is deterministic
// but arbitrary.
func (g *Graph) BuildList() []MvsVersion {
	seenRoot := make(map[string]bool, len(g.roots))

	var list []MvsVersion
	for _, r := range g.roots {
		if seenRoot[r.PackageName] {
			// Multiple copies of the same root, with the same or different versions,
			// are a bit of a degenerate case: we will take the transitive
			// requirements of both roots into account, but only the higher one can
			// possibly be selected. However — especially given that we need the
			// seenRoot map for later anyway — it is simpler to support this
			// degenerate case than to forbid it.
			continue
		}

		if v := g.Selected(r.PackageName); !v.IsNone() {
			list = append(list, MvsVersion{PackageName: r.PackageName, PackageVersion: v})
		}
		seenRoot[r.PackageName] = true
	}
	uniqueRoots := list

	for path, version := range g.selected {
		if !seenRoot[path] {
			list = append(list, MvsVersion{PackageName: path, PackageVersion: version})
		}
	}
	ModSort(g.cmp, list[len(uniqueRoots):])

	return list
}

// WalkBreadthFirst invokes f once, in breadth-first order, for each module
// version other than "none" that appears in the graph, regardless of whether
// that version is selected.
func (g *Graph) WalkBreadthFirst(f func(m MvsVersion)) {
	var queue []MvsVersion
	enqueued := make(map[MvsVersion]bool)
	for _, m := range g.roots {
		if m.Version != "none" {
			queue = append(queue, m)
			enqueued[m] = true
		}
	}

	for len(queue) > 0 {
		m := queue[0]
		queue = queue[1:]

		f(m)

		reqs, _ := g.RequiredBy(m)
		for _, r := range reqs {
			if !enqueued[r] && r.Version != "none" {
				queue = append(queue, r)
				enqueued[r] = true
			}
		}
	}
}

// FindPath reports a shortest requirement path starting at one of the roots of
// the graph and ending at a module version m for which f(m) returns true, or
// nil if no such path exists.
func (g *Graph) FindPath(f func(MvsVersion) bool) []MvsVersion {
	// firstRequires[a] = b means that in a breadth-first traversal of the
	// requirement graph, the module version a was first required by b.
	firstRequires := make(map[MvsVersion]MvsVersion)

	queue := g.roots
	for _, m := range g.roots {
		firstRequires[m] = MvsVersion{}
	}

	for len(queue) > 0 {
		m := queue[0]
		queue = queue[1:]

		if f(m) {
			// Construct the path reversed (because we're starting from the far
			// endpoint), then reverse it.
			path := []MvsVersion{m}
			for {
				m = firstRequires[m]
				if m.PackageName == "" {
					break
				}
				path = append(path, m)
			}

			i, j := 0, len(path)-1
			for i < j {
				path[i], path[j] = path[j], path[i]
				i++
				j--
			}

			return path
		}

		reqs, _ := g.RequiredBy(m)
		for _, r := range reqs {
			if _, seen := firstRequires[r]; !seen {
				queue = append(queue, r)
				firstRequires[r] = m
			}
		}
	}

	return nil
}
