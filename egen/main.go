package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ostlerc/voter/election"
)

var (
	peak       = flag.Bool("peak", false, "Force generation with peak preference")
	cond       = flag.Bool("cond", false, "Force condorcet winner")
	pref       = flag.Bool("pref", false, "Force preference of some candidate")
	weight     = flag.Int("weight", 5, "Maximum weight of a vote")
	Votes      = flag.Int("vote", 4, "Number of voters in election")
	Candidates = flag.Int("cand", 3, "Number of candidates in election")
	random     = flag.Bool("rand", false, "Use a random vote/cand count")
	o          = flag.String("o", "json", "output type [json,csv]")

	voter = Voter(&RandVoter{})
)

const (
	i = "csv" //only type so hard coded not a flag for now
)

func init() {
	flag.Parse()
	election.Setup()

	if *random {
		*Votes = election.Max(election.R.Intn(*Votes), 3)           //minimum of 3 votes
		*Candidates = election.Max(election.R.Intn(*Candidates), 3) //minimum of 3 candidates
	}
}

func Campaign() *election.Election {
	e := election.New(*Votes, *Candidates)
	if *pref {
		e.F = election.NewPref(*Candidates)
		voter = &PreferVoter{
			f:    e.F,
			peak: *peak,
		}
	} else if *peak {
		voter = &PeakVoter{}
	}

	return e
}

func Fix(e *election.Election) {
	if *pref && e.F == nil {
		for _, v := range e.V {
			v.Prefer(e.F)
		}
	}
	c := e.Condorcet()

	for *cond && c == -1 { //continually make more until we have a condorcet winner. slow but it works
		i := election.R.Intn(*Votes)
		e.V[i] = voter.Vote(*Candidates)
		c = e.Condorcet()
	}

	e.F = e.Pref()

	if c != -1 {
		e.C = &c
	}

	p := e.Peak()
	if p != -1 {
		e.P = &p
	}
}

func main() {
	if *Candidates < 2 {
		log.Fatal("Invalid number of candidates, must be >1")
	}
	if *Votes < 2 {
		log.Fatal("Invalid number of votes, must be >1 ")
	}

	if *o != "json" && *o != "csv" {
		log.Fatal("Invalid output type '", *o, "'")
	}
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	e := Campaign()

	if stat.Mode()&os.ModeCharDevice == 0 { //read in csv and create election from it
		if *peak {
			log.Fatal("Incorrect usage - cannot force peak with given data set, randomly create an election with force peak instead")
		}
		e = election.ParseFrom(i, os.Stdin)
		*Votes = len(e.V)
		*Candidates = e.N
		e.F = e.Pref()
	} else {
		for i := 0; i < *Votes; i++ {
			e.V[i] = voter.Vote(*Candidates)
		}
	}

	Fix(e)

	if *o == "json" {
		dat, err := json.Marshal(&e)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dat))
	} else if *o == "csv" {
		fmt.Println(e.CSV())
	}
}
