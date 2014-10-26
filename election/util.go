package election

import "strconv"

type IntPair struct {
	First  int `json:"a"`
	Second int `json:"b"`
}

func Strn(n int) string {
	return strconv.Itoa(R.Intn(n))
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
