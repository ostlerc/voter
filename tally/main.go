package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ostlerc/voter/election"
)

var (
	t = flag.String("t", csv(election.TallyKeys()), "tally type results")
	o = flag.String("o", "json", "output type [json,csv]")
)

func csv(l []string) string {
	res := ""
	for _, s := range l {
		res += s + ","
	}
	return res[:len(res)-1]
}

type TallyResult struct {
	Results map[string][]int  `json:"results"`
	Names   map[string]string `json:"names,omitempty"`
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

	reader := bufio.NewReader(os.Stdin)
	dat, err := reader.ReadBytes('\n')
	if err != nil { //assume we don't generate elections outside of buffer range
		log.Fatal(err)
	}
	e := &election.Election{}
	err = json.Unmarshal(dat, e)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string][]int)
	for _, t := range talliers {
		m[t.Key()] = t.Tally(e)
	}

	if *o == "json" {
		dat, err = json.Marshal(&TallyResult{Results: m, Names: e.M})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dat))
	} else {
		fmt.Print("rank,")
		for i := 0; i < len(talliers); i++ {
			fmt.Print(talliers[i].Key())
			if i+1 != len(talliers) {
				fmt.Print(",")
			}
		}
		fmt.Println("")
		end := len(m[talliers[0].Key()])
		keys := make([]string, 0) //making this list guarantees ordering. Range on map has no guaranteed order
		for k, _ := range m {
			keys = append(keys, k)
		}
		for i := 0; i < end; i++ {
			for count, k := range keys {
				if count == 0 {
					fmt.Print(i+1, ",")
				}
				if len(e.M) > 0 {
					fmt.Print(e.M[strconv.Itoa(m[k][i])])
				} else {
					fmt.Print(m[k][i])
				}
				if count+1 != len(m) {
					fmt.Print(",")
				}
			}
			fmt.Println("")
		}
	}
}
