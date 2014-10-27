package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ostlerc/voter/election"
)

var (
	o = flag.String("o", "dot", "graph output type. [json,dot]")
	i = flag.String("i", "json", "graph input type. ["+election.CSVFlat(election.Parsers())+"]")
)

func main() {
	flag.Parse()
	election.Setup()

	if *o != "dot" && *o != "json" {
		log.Fatal("invalid type ", *o)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	hasStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if !hasStdin {
		log.Fatal("No stdin to read")
	}

	e := election.ParseFrom(*i, os.Stdin)
	if e == nil {
		log.Fatal("Incorrect input type '", *i, "'")
	}

	edges := e.Graph().Edges()
	if *o == "dot" {
		fmt.Println(edges.Dot())
	} else if *o == "json" {
		fmt.Println(edges.JSON())
	}
}
