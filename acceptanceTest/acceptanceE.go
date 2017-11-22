package acceptanceTest

// New returns an error that formats as the given text.
func New(text string) error {
	return &acceptanceE{text}
}

// errorString is a trivial implementation of error.
type acceptanceE struct {
	s string
}

func (e *acceptanceE) Error() string {
	return e.s
}
