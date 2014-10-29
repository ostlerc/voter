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
	o = flag.String("o", "csv", "output type [json,csv]")
	i = flag.String("i", "json", "tally input type. ["+election.CSVFlat(election.Parsers())+"]")
	v = flag.Bool("v", false, "verbose json output. Shows election and manipulation values")
)

type TallyResult struct {
	Results     map[string][]int                  `json:"results"`
	Names       map[string]string                 `json:"names,omitempty"`
	Condorcet   *int                              `json:"condorcet,omitempty"`
	Election    *election.Election                `json:"election,omitempty"`
	M           map[string]*election.Manipulation `json:"manipulations,omitempty"`
	CondoretWon bool                              `json:"condorcet_won"`
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

	manipulations := make(map[string]*election.Manipulation)

	if *v {
		for _, t := range talliers {
			m := e.FindManipulation(t)
			if m != nil {
				manipulations[t.Key()] = m
			}
		}
	}

	if *o == "json" {
		res := &TallyResult{
			Results: m,
			Names:   e.M,
		}

		if c := e.Condorcet() + 1; c != 0 {
			res.Condorcet = &c
		}

		if *v {
			res.M = manipulations
			res.Election = e
		}

		dat, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dat))
	} else {
		if *v {
			fmt.Println(e.CSV())
			fmt.Println("")
		}

		keys := make(sort.StringSlice, 0) //making this list guarantees ordering. Range on map has no guaranteed order
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Sort(keys)

		fmt.Print("rank,")
		for i := 0; i < len(keys); i++ {
			fmt.Print(keys[i])
			if i+1 != len(talliers) {
				fmt.Print(",")
			}
		}
		fmt.Println("")

		if *v {
			fmt.Print("manipulation,")
			for i := 0; i < len(keys); i++ {
				_, ok := manipulations[keys[i]]
				fmt.Print(ok, ",")
				if i+1 != len(talliers) {
					fmt.Print(",")
				}
			}
			fmt.Println("")

			if e.C != nil {
				fmt.Print("condorcet won,")
				for i := 0; i < len(keys); i++ {
					fmt.Print(m[keys[i]][0] == *e.C)
					if i+1 != len(talliers) {
						fmt.Print(",")
					}
				}
				fmt.Println("")
			}
		}
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
					if m[k][i] == -1 {
						fmt.Print(m[k][i])
					} else {
						fmt.Print(m[k][i] + 1)
					}
				}
				if j+1 != len(m) {
					fmt.Print(",")
				}
			}
			fmt.Println("")
		}

		if f := e.Pref(); *v && f != nil {
			continuedPref := true

			for j := 0; j < len(keys) && continuedPref; j++ {
				for _, v := range m[keys[j]] {
					if v == f.First {
						break
					}
					if v == f.Second {
						continuedPref = false
						break
					}
					if keys[j] == "bucklin" {
						break
					}
				}
			}

			fmt.Printf("pref intact,%v\n", continuedPref)
		}

	}
}
