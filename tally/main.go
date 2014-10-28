package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ostlerc/voter/election"
)

var (
	t = flag.String("t", election.CSVFlat(election.TallyKeys()), "tally type results")
	o = flag.String("o", "json", "output type [json,csv]")
	i = flag.String("i", "json", "tally input type. ["+election.CSVFlat(election.Parsers())+"]")
)

type TallyResult struct {
	Results map[string][]int  `json:"results"`
	Names   map[string]string `json:"names,omitempty"`
}

func main() {
	flag.Parse()
	election.Setup()

	if *o != "json" && *o != "csv" {
		log.Fatal("Invalid output type '", *o, "'")
	}
	var talliers []election.Tallier
	for _, key := range strings.Split(*t, ",") {
		if tally := election.GetTally(key); tally == nil {
			log.Fatal("invalid tally type '", *t, "'")
		} else {
			talliers = append(talliers, tally)
		}
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if !(stat.Mode()&os.ModeCharDevice == 0) {
		log.Fatal("No stdin to read. Expecting json election as stdin.")
	}

	e := election.ParseFrom(*i, os.Stdin)

	m := make(map[string][]int)
	for _, t := range talliers {
		m[t.Key()] = t.Tally(e)
	}

	if *o == "json" {
		dat, err := json.Marshal(&TallyResult{Results: m, Names: e.M})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dat))
	} else {
		fmt.Print("rank,")
		keys := make(sort.StringSlice, 0) //making this list guarantees ordering. Range on map has no guaranteed order
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Sort(keys)
		for i := 0; i < len(keys); i++ {
			fmt.Print(keys[i])
			if i+1 != len(talliers) {
				fmt.Print(",")
			}
		}
		fmt.Println("")
		end := len(m[talliers[0].Key()])
		for i := 0; i < end; i++ {
			for j := 0; j < len(keys); j++ {
				if j == 0 {
					fmt.Print(i+1, ",")
				}
				k := keys[j]
				if len(e.M) > 0 {
					name := e.M[strconv.Itoa(m[k][i])]
					if name == "" {
						fmt.Print(m[k][i])
					} else {
						fmt.Print(name)
					}
				} else {
					fmt.Print(m[k][i])
				}
				if j+1 != len(m) {
					fmt.Print(",")
				}
			}
			fmt.Println("")
		}
	}
}
