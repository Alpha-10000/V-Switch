package main

import "flag"
import "fmt"
import "os"

const minArgs = 2
type Opts struct {
	
}

func parseCmd(opts *Opts)  {
	flag.Parse()

	if flag.NArg() < minArgs {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
}

