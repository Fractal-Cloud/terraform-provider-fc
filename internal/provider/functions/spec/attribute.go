package spec

// Attribute is one attribute of a component function's config, mapped to
// the component parameter the catalog declares for it.
type Attribute struct {
	// Name is the attribute's name in Terraform.
	Name string
	// Key is the parameter key the attribute is written under.
	Key string
	// Aliases are further keys the value is also written under, for offers
	// that spell the same setting differently.
	Aliases []string
	Kind    Kind
	// Required makes the attribute mandatory in the type constraint.
	Required bool
	// OneOf, when set, lists the only values a String attribute accepts.
	OneOf []string
	// Fields are the attributes of an Object attribute. Nested fields are
	// always optional.
	Fields []Attribute
}
