package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ostlerc/voter/election"
)

var (
	o = flag.String("o", "dot", "graph output type. [json,dot]")
	i = flag.String("i", "json", "graph input type. [json,csv]")
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

	e := election.New(0, 0)
	if *i == "json" {
		reader := bufio.NewReader(os.Stdin)
		dat, err := reader.ReadBytes('\n')
		if err != nil { //assume we don't generate elections outside of buffer range
			log.Fatal(err)
		}
		err = json.Unmarshal(dat, e)
		if err != nil {
			log.Fatal(err)
		}
	} else if *i == "csv" {
		election.CSVElection(e, os.Stdin)
	} else {
		log.Fatal("Incorrect input type '", *i, "'")
	}
	edges := e.Graph().Edges()
	if *o == "dot" {
		fmt.Println(edges.Dot())
	} else if *o == "json" {
		fmt.Println(edges.JSON())
	}
}
