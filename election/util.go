package election

import (
	"encoding/json"
	"strconv"
)

type Ints []int
type VoteAr []*Vote
type VoteM []VoteAr

type IntPair struct {
	First  int `json:"a"`
	Second int `json:"b"`
}

func Strn(n int) string {
	return strconv.Itoa(R.Intn(n))
}

func ArEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Contains(v int, l []int) bool {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return true
		}
	}
	return false
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RemoveAt(i int, l []int) []int {
	return append(l[:i], l[i+1:]...)
}

func RemoveValue(v int, l []int) []int {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return RemoveAt(i, l)
		}
	}
	panic("No value")
}

func Index(v int, l []int) int {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return i
		}
	}
	return -1
}

func allPermutations(l []int) [][]int {
	if len(l) == 1 {
		return [][]int{l}
	}

	res := make([][]int, 0)
	for i := 0; i < len(l); i++ {
		next := make([]int, len(l))            //allocate new memory
		copy(next, l)                          //copy l
		next = append(next[:i], next[i+1:]...) //remove ith element
		perms := allPermutations(next)
		for j, p := range perms {
			perms[j] = append([]int{l[i]}, p...)
		}
		res = append(res, perms...)
	}
	return res
}

// Perms returns all permutations for a number of candidates
// This is mostly just a wrapper around the allPermutations internal function
func Perms(n int) [][]int {
	set := make([]int, n)
	for i := 0; i < n; i++ {
		set[i] = i
	}

	return allPermutations(set)
}

func CSVFlat(l []string) string {
	res := ""
	for _, s := range l {
		res += s + ","
	}
	return res[:len(res)-1]
}

func (i Ints) Max() int {
	if len(i) == 0 {
		panic("Max on empty list")
	}
	m := i[0]
	for _, v := range i {
		if v > m {
			m = v
		}
	}
	return m
}

func (i Ints) Min() int {
	if len(i) == 0 {
		panic("Max on empty list")
	}
	m := i[0]
	for _, v := range i {
		if v < m {
			m = v
		}
	}
	return m
}

func (i Ints) Maxi(ig Ints) int {
	if len(i) == 0 {
		panic("Max on empty list")
	}
	m := -1
	for k, v := range i {
		if Contains(k, ig) {
			continue
		}
		if m == -1 || v > i[m] {
			m = k
		}
	}
	return m
}

func (i Ints) Mini(ig Ints) int {
	if len(i) == 0 {
		panic("Max on empty list")
	}
	m := -1
	for k, v := range i {
		if Contains(k, ig) {
			continue
		}
		if m == -1 || v < i[m] {
			m = k
		}
	}
	return m
}

func (i Ints) Minus(ignore Ints) Ints {
	res := make(Ints, len(i))
	copy(res, i)
	for _, v := range i {
		if Contains(v, i) {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (v VoteM) Maxi(ignore Ints) int {
	max := -1
	maxi := -1
	for i, v := range v {
		if Contains(i, ignore) {
			continue
		}
		if max == -1 || max < len(v) {
			maxi = i
			max = len(v)
		}
	}
	return maxi
}

func (v VoteM) Mini(ignore Ints) int {
	min := -1
	mini := -1
	for i, v := range v {
		if Contains(i, ignore) {
			continue
		}
		if min == -1 || min >= len(v) {
			mini = i
			min = len(v)
		}
	}
	return mini
}

func (v VoteM) Lens() []int {
	res := make([]int, len(v))
	for i, v := range v {
		res[i] = len(v)
	}
	return res
}

func tojson(i interface{}) string {
	dat, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func GetName(M map[string]string, v string) string {
	if s, ok := M[v]; ok {
		return s
	}
	return v
}
