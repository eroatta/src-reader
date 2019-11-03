package expander

type expander struct {
	name string
}

// Name returns the name of the Expander.
func (e expander) Name() string {
	return e.name
}
