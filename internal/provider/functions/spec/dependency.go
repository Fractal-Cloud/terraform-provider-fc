package spec

// Dependency is a component-reference attribute whose target the component
// cannot exist without.
type Dependency struct {
	// Name is the attribute's name in Terraform.
	Name string
	// ComponentType is the type the referenced component must have.
	ComponentType string
}
