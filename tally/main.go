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
	v = flag.Bool("v", false, "verbose output. Show all tally information")
)

type IrrelevantCand struct {
	Candidate     string `json:"candidate"`
	ChangesWinner bool   `json:"causes_change"`
}

type TallyResult struct {
	Results     map[string][]int                  `json:"results"`
	Names       map[string]string                 `json:"names,omitempty"`
	Condorcet   *int                              `json:"condorcet,omitempty"`
	Election    *election.Election                `json:"election,omitempty"`
	M           map[string]*election.Manipulation `json:"manipulations,omitempty"`
	CondoretWon bool                              `json:"condorcet_won"`
	PrefIntact  *bool                             `json:"pref_intact,omitempty"`
	Irrelevant  *IrrelevantCand                   `json:"irr_cand,omitempty"`
}

func irrelevant(e *election.Election, tallies map[string][]int) *IrrelevantCand {
	found := make(map[int]bool)
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
		res := t.Tally(e2)
		if v[0] != res[0] {

			return &IrrelevantCand{
				Candidate:     name,
				ChangesWinner: true,
			}
		}
	}

	return &IrrelevantCand{
		Candidate:     name,
		ChangesWinner: false,
	}
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

	prefIntact := true
	if e.F == nil {
		e.F = e.Pref()
	}

	irrCand := irrelevant(e, m)

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
			for j := 0; j < len(keys) && prefIntact; j++ {
				for _, v := range m[keys[j]] {
					if v == e.F.First {
						break
					}
					if v == e.F.Second {
						prefIntact = false
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
			if e.F != nil {
				res.PrefIntact = &prefIntact
			}
			res.Irrelevant = irrCand
		}

		dat, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dat))
	} else {
		if *v {
			fmt.Println(e.CSV())
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

		if *v {
			if e.F != nil {
				fmt.Printf("pref intact,%v\n", prefIntact)
			}
			if irrCand != nil {
				fmt.Printf("irrlvnt alters,%v,%v\n", irrCand.Candidate, irrCand.ChangesWinner)
			}
		}
	}
}
