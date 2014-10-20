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

	orderer = Orderer(&RandOrder{})
)

func init() {
	flag.Parse()
}

func main() {
	e := election.New()
	if *pref {
		f := &election.Favorite{ //prefer a over b (always)
			A: election.R.Intn(*election.Candidates),
			B: election.R.Intn(*election.Candidates),
		}

		for f.A == f.B { //verify unique random values
			f.B = election.R.Intn(*election.Candidates)
		}

		e.F = f
		orderer = &PreferOrder{
			f:    f,
			peak: *peak,
		}
	}

	fmt.Println(`{"peak":"`, *peak, `"}`)
	e.P = *peak

	for i := 0; i < *election.Votes; i++ {
		v := election.NewVote()
		o := orderer.Order(*election.Candidates)
		for j := 0; j < len(o); j++ {
			v.C[strconv.Itoa(j)] = o[j]
		}
		e.V[i] = v
	}

	dat, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dat))
}
