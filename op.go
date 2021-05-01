package myopa

type Op int

func (op Op) String() string {
	return [...]string{
		"invalid",
		"equal",
	}[op]
}

const (
	OpEq Op = iota + 1
)
