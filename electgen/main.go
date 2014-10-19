package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"github.com/ostlerc/voter/election"
)

var (
	peak = flag.Bool("peak", false, "Force generation with peak preference")
	cond = flag.Bool("cond", false, "Force condorcet winner")
	pref = flag.Bool("pref", false, "Force preference of some candidate")
)

func init() {
	flag.Parse()
}

func main() {
	e := election.New()
	for i := 0; i < *election.Votes; i++ {
		v := election.NewVote()
		candidates := make([]int, *election.Candidates)
		for j := 0; j < *election.Candidates; j++ {
			candidates[j] = j
		}
		for j := 0; j < *election.Candidates; j++ {
			k := election.RandCand(candidates)
			v.C[strconv.Itoa(j)] = candidates[k]
			candidates = election.RemoveAt(k, candidates)
		}
		e.V[i] = v
	}

	dat, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dat))
}
