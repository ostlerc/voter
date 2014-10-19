package main

// Picker picks candidates
type Picker interface {
	Pick([]int) int
}
