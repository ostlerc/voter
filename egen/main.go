package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/ostlerc/voter/election"
)

var (
	peak       = flag.Bool("peak", false, "Force generation with peak preference")
	cond       = flag.Bool("cond", false, "Force condorcet winner")
	pref       = flag.Bool("pref", false, "Force preference of some candidate")
	Votes      = flag.Int("vote", 6, "Number of voters in election")
	Candidates = flag.Int("cand", 3, "Number of candidates in election")
	random     = flag.Bool("rand", false, "Use a random vote/cand count")

	voter = Voter(&RandVoter{})
)

func init() {
	flag.Parse()
	election.Setup()

	if *random {
		*Votes = election.Max(election.R.Intn(*Votes), 3)           //minimum of 3 votes
		*Candidates = election.Max(election.R.Intn(*Candidates), 3) //minimum of 3 candidates
	}
}

func main() {
	e := election.New(*Votes, *Candidates)
	if *pref {
		f := &election.IntPair{ //prefer a over b (always)
			First:  election.R.Intn(*Candidates),
			Second: election.R.Intn(*Candidates),
		}

		for f.First == f.Second { //verify unique random values
			f.Second = election.R.Intn(*Candidates)
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

	for i := 0; i < *Votes; i++ {
		e.V[i] = voter.Vote(*Candidates)
	}

	c, ranks := e.Condorcet()

	for *cond && c == -1 { //continually make more until we have a condorcet winner. slow but it works
		i := election.R.Intn(*Votes)
		e.V[i] = voter.Vote(*Candidates)
		c, ranks = e.Condorcet()
	}

	if c != -1 {
		e.C = &c
	}
	e.R = ranks

	dat, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dat))
}
