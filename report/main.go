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
	o = flag.String("o", "dot", "report output type. [json,dot]")
)

func main() {
	flag.Parse()
	election.Setup()

	if *o != "dot" && *o != "json" {
		log.Fatal("invalid type", *o)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	hasStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if !hasStdin {
		log.Fatal("No stdin to read")
	}

	reader := bufio.NewReader(os.Stdin)
	dat, err := reader.ReadBytes('\n')
	if err != nil { //assume we don't generate elections outside of buffer range
		log.Fatal(err)
	}
	e := &election.Election{}
	err = json.Unmarshal(dat, e)
	if err != nil {
		log.Fatal(err)
	}
	edges := e.Graph().Edges()
	if *o == "dot" {
		fmt.Println(edges.Dot())
	} else if *o == "json" {
		fmt.Println(edges.JSON())
	}
}
