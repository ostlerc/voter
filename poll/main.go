package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ostlerc/voter/election"
)

type TallySum struct {
	PrefChanged int            `json:"pref_changed"`
	M           map[string]int `json:"manipulations"`
	Irrelevant  int            `json:"irr_cand_affect"`
	Total       int            `json:"total"`
}

func main() {
	r := bufio.NewReader(os.Stdin)
	dat, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	var res election.TallyResult
	sum := &TallySum{
		M: make(map[string]int),
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

		if res.PrefIntact != nil && !*res.PrefIntact {
			sum.PrefChanged++
		}

		if res.Irrelevant != nil && res.Irrelevant.ChangesWinner {
			sum.Irrelevant++
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
