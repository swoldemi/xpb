package main

import (
	"github.com/swoldemi/xpb/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
