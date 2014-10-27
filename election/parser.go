package election

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	parsers = make(map[string]ElectionParser)
)

type ElectionParser interface {
	Parse(io.Reader) *Election
}

type csvParser struct{}
type jsonParser struct{}

func init() {
	RegisterParser("csv", &csvParser{})
	RegisterParser("json", &jsonParser{})
}

func Parsers() []string {
	res := make([]string, 0)
	for k, _ := range parsers {
		res = append(res, k)
	}
	return res
}

func RegisterParser(key string, p ElectionParser) {
	parsers[key] = p
}

func ParseFrom(key string, r io.Reader) *Election {
	if p, ok := parsers[key]; ok {
		return p.Parse(r)
	} else {
		return nil
	}
}

func (*csvParser) Parse(r io.Reader) *Election {
	e := New(0, 0)
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	file := strings.Split(string(dat), "\n")

	cand := len(file) - 1
	e.V = make([]*Vote, 0)

	for y, row := range file {
		csv := strings.Split(row, ",")
		if len(csv) < 2 {
			cand--
			continue
		}
		for x, col := range csv {
			if y == 0 { //handle weights
				if x == 0 {
					continue //empty anyways
				}
				iw, err := strconv.Atoi(col)
				if err != nil {
					log.Fatal(err)
				}
				e.V = append(e.V, &Vote{C: make(map[string]int), W: iw})
				continue
			}
			if x == 0 { //handle name map
				e.M[strconv.Itoa(y-1)] = col
				continue
			}

			icol, err := strconv.Atoi(col)
			if err != nil {
				log.Fatal(err)
			}
			e.V[x-1].C[strconv.Itoa(y-1)] = icol - 1

		}
	}

	e.N = cand

	return e
}

func (*jsonParser) Parse(r io.Reader) *Election {
	reader := bufio.NewReader(r)
	dat, err := reader.ReadBytes('\n')
	if err != nil { //assume we don't generate elections outside of buffer range
		log.Fatal(err)
	}
	e := New(0, 0)
	err = json.Unmarshal(dat, e)
	if err != nil {
		log.Fatal(err)
	}
	return e
}
