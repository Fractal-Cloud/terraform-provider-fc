// Package spec builds component functions from a declarative description of
// the catalog entry they produce, for components whose functions differ only
// in their attributes.
package spec

// Spec describes a component function.
type Spec struct {
	// Name is the function name (provider::fc::<Name>).
	Name string
	// ComponentType is the catalog type the function builds.
	ComponentType string
	Summary       string
	Description   string
	Attributes    []Attribute
	Dependencies  []Dependency
	// Links adds a links attribute.
	Links bool
}
