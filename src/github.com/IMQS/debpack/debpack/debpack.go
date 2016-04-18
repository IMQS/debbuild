package main

import (
	"fmt"
	"github.com/IMQS/debpack/pack"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage debpack <deffile.json>")
		os.Exit(0)
	}
	deb, err := pack.NewDebBuild(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}

	err = deb.Build()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
}
