package dockerDriver

func MockSearchContainer(name string) (string, bool) {
	return RandSeq(10), true
}
