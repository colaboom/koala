package main

type DirGenerator struct {
}

func (dg *DirGenerator) Run(opt *Option) (err error) {
	return
}

func init() {
	dir := &DirGenerator{}

	Register("dir generator", dir)
}
