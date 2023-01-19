package main

import (
	"github.com/antlabs/cat"
	"github.com/guonaihong/clop"
)

func main() {
	c := cat.Cat{}
	clop.Bind(&c)
	c.Main()
}
