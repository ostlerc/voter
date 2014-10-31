package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ostlerc/voter/election"
)

type TallySum struct {
	PrefChanged   map[string]int `json:"pref_changed,omitempty"`
	M             map[string]int `json:"manipulations,omitempty"`
	Irrelevant    map[string]int `json:"irr_cand_affect,omitempty"`
	CondorcetLost map[string]int `json:"condorcet_lost,omitempty"`
	Total         int            `json:"total"`
}

func main() {
	flag.Parse()
	r := bufio.NewReader(os.Stdin)
	dat, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	var res election.TallyResult
	sum := &TallySum{
		M:             make(map[string]int),
		PrefChanged:   make(map[string]int),
		Irrelevant:    make(map[string]int),
		CondorcetLost: make(map[string]int),
	}
	for err == nil && len(dat) > 0 {
		err := json.Unmarshal(dat, &res)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range res.M {
			if v != nil {
				sum.M[k]++
			}
		}

		for k, v := range res.PrefIntact {
			if !v {
				sum.PrefChanged[k]++
			}
		}

		for k, v := range res.Irrelevant {
			if v.ChangesWinner {
				sum.Irrelevant[k]++
			}
		}

		if res.Condorcet != nil {
			for k, v := range res.Results {
				if v[0]+1 != *res.Condorcet && v[0] != -1 {
					sum.CondorcetLost[k]++
				}
			}
		}

		sum.Total++

		dat, err = r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
	}

	dat, err = json.Marshal(&sum)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(dat))
}
