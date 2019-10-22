package splitter

type splitter struct {
	name string
}

// Name returns the name of the splitter.
func (s splitter) Name() string {
	return s.name
}
