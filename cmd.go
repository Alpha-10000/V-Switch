package main

import "flag"
import "fmt"
import "os"

type Opts struct {
	intfs []string
}

const minArgs = 2

func parseCmd(opts *Opts)  {
	flag.Parse()

	if flag.NArg() < minArgs {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	opts.intfs = make([]string, 0)

	for i := 0; i < flag.NArg(); i++ {
		opts.intfs = append(opts.intfs, flag.Arg(i))
	}
}
