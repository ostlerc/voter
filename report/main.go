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
	t = flag.String("type", "dot", "Default type. Options: [json, dot]")
)

func main() {
	flag.Parse()

	if *t != "dot" && *t != "json" {
		log.Fatal("invalid type", *t)
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
	if *t == "dot" {
		fmt.Println(edges.Dot())
	} else if *t == "json" {
		fmt.Println(edges.JSON())
	}
}
