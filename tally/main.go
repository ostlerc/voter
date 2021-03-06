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
	t    = flag.String("t", election.CSVFlat(election.TallyKeys()), "tally type results")
	o    = flag.String("o", "csv", "output type [json,csv]")
	i    = flag.String("i", "json", "tally input type. ["+election.CSVFlat(election.Parsers())+"]")
	v    = flag.Bool("v", false, "verbose output. Show all tally information")
	dumb = flag.Bool("dumb", false, "Use dumb scoring (only check if manipulation changes top candidate to your top candidate")
)

func irrelevant(e *election.Election, tallies map[string][]int) map[string]*election.IrrelevantCand {
	found := make(map[int]bool)
	res := make(map[string]*election.IrrelevantCand)
	for _, v := range tallies {
		found[v[0]] = true //not in a first place, then you are IRRELEVANT!
	}
	cand := -1
	for i := 0; i < e.N; i++ {
		if _, ok := found[i]; !ok {
			cand = i
			break
		}
	}
	if cand == -1 {
		return nil
	}

	name := strconv.Itoa(cand)

	if len(e.M) > 0 {
		name = e.M[name]
	}

	e2 := e.RemoveCandidate(cand)
	for k, v := range tallies {
		t := election.GetTally(k)
		tres := t.Tally(e2)
		if v[0] != tres[0] {
			res[k] = &election.IrrelevantCand{
				Candidate:     name,
				ChangesWinner: true,
			}
		}
	}

	return res
}

func main() {
	flag.Parse()
	election.Setup()

	election.DumbScore = *dumb

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

	prefIntact := make(map[string]bool)
	if e.F == nil {
		e.F = e.Pref()
	}

	irrCands := irrelevant(e, m)

	keys := make(sort.StringSlice, 0) //making this list guarantees ordering. Range on map has no guaranteed order
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	if *v {
		for _, t := range talliers {
			m := e.FindManipulation(t)
			if m != nil {
				manipulations[t.Key()] = m
			}
		}

		if e.F != nil {
			for j := 0; j < len(keys); j++ {
				prefIntact[keys[j]] = true
				for _, v := range m[keys[j]] {
					if v == e.F.First {
						break
					}
					if v == e.F.Second {
						prefIntact[keys[j]] = false
						break
					}
					if keys[j] == "bucklin" {
						break
					}
				}
			}
		}
	}

	if *o == "json" {
		res := &election.TallyResult{
			Results:    m,
			Names:      e.M,
			PrefIntact: make(map[string]bool),
			Irrelevant: make(map[string]*election.IrrelevantCand),
		}

		if c := e.Condorcet() + 1; c != 0 {
			res.Condorcet = &c
		}

		if *v {
			res.M = manipulations
			res.Election = e
			if e.F != nil {
				for k, v := range prefIntact {
					res.PrefIntact[k] = v
				}
			}
			for k, v := range irrCands {
				res.Irrelevant[k] = v
			}
		}

		dat, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dat))
	} else {
		if *v {
			fmt.Println(e.CSV())

			f := e.Pref()
			if f != nil {
				fmt.Printf("\npref,%d,%d", f.First+1, f.Second+1)
			}
			c := e.Condorcet()
			if c != -1 {
				fmt.Printf("\ncondorcet,%d\n", c+1)
			}
		}

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
				fmt.Print(ok)
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

			if e.F != nil {
				fmt.Printf("pref intact,")
				for j := 0; j < len(keys); j++ {
					if v, ok := prefIntact[keys[j]]; ok && v {
						fmt.Print("true")
					} else {
						fmt.Print("false")
					}
					if j+1 != len(prefIntact) {
						fmt.Print(",")
					}
				}
				fmt.Println("")
			}
			fmt.Printf("irrlvnt alters,")
			for j := 0; j < len(keys); j++ {
				if _, ok := irrCands[keys[j]]; ok {
					fmt.Print("true")
				} else {
					fmt.Print("false")
				}
				if j+1 != len(irrCands) {
					fmt.Print(",")
				}

			}
			fmt.Println("")
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
	}
}
