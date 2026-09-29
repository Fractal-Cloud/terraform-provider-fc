package spec

// Kind is the Terraform type of an attribute and how its value is written
// as a component parameter.
type Kind int

const (
	// String is written as is.
	String Kind = iota
	// Int64 is written as its decimal text.
	Int64
	// Bool is written as "true" or "false".
	Bool
	// StringList is written as a JSON array.
	StringList
	// StringMap is written as a JSON object.
	StringMap
	// Object is a nested object of Fields, written as a JSON object keyed by
	// the fields' parameter keys.
	Object
)
