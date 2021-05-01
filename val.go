package myopa

// Val is a value
type Val struct {
	T VT
	V interface{}
}

// VT is a value type
type VT int

func (vt VT) String() string {
	return [...]string{
		"invalid",
		"constant",
		"key-value",
	}[vt]
}

// List of value types
const (
	VTConstant VT = iota + 1
	VTKeyValue
)
