package miner

// New TODO
func New(name string) interface{} {
	var miner interface{}
	switch name {
	case "samurai":
		miner = NewSamuraiExtractor()
	}

	return miner
}
