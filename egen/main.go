package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/ostlerc/voter/election"
)

var (
	peak = flag.Bool("peak", false, "Force generation with peak preference")
	cond = flag.Bool("cond", false, "Force condorcet winner")
	pref = flag.Bool("pref", false, "Force preference of some candidate")

	voter = Voter(&RandVoter{})
)

func init() {
	flag.Parse()
	election.Setup()
}

func main() {
	e := election.New()
	if *pref {
		f := &election.IntPair{ //prefer a over b (always)
			First:  election.R.Intn(*election.Candidates),
			Second: election.R.Intn(*election.Candidates),
		}

		for f.First == f.Second { //verify unique random values
			f.Second = election.R.Intn(*election.Candidates)
		}

		e.F = f
		voter = &PreferVoter{
			f:    f,
			peak: *peak,
		}
	} else if *peak {
		voter = &PeakVoter{}
	}

	if *peak {
		e.P = peak
	}

	for i := 0; i < *election.Votes; i++ {
		e.V[i] = voter.Vote(*election.Candidates)
	}

	c := e.Condorcet()
	if c != -1 {
		e.C = &c
	}
	e.R = e.Rank()

	dat, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dat))
}
