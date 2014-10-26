package election

func (e *Egraph) Slater() ([]int, []int) {
	var minE, minW int
	var minEdges []int
	var minWeights []int
	for i, p := range Perms(len(e.Nodes)) {
		d, w := e.slater(p)
		if i == 0 {
			minEdges = p
			minWeights = p
			minE, minW = d, w
			continue
		}
		if d < minE {
			minEdges = p
			minE = d
		}
		if w < minW {
			minWeights = p
			minW = w
		}

		if minE == 0 && minW == 0 { //can't get better than that
			break
		}
	}

	return minEdges, minWeights
}

//Returns slater value (disagreeing edges, disagreement weights) of given permutation
func (e *Egraph) slater(p []int) (int, int) {
	edges := 0
	weights := 0

	for i := 0; i < len(p); i++ {
		for j := i + 1; j < len(p); j++ {
			v := e.Nodes[p[j]][p[i]]
			if v > 0 {
				edges++
				weights += v
			}
		}
	}

	return edges, weights
}
