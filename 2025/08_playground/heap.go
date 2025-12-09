package main

type pairHeap []pair

func (h *pairHeap) len() int {
	return len(*h)
}

func (h *pairHeap) push(p pair) {
	// add new element at the end
	*h = append(*h, p)
	// percolate up while larger than parent
	i := h.len() - 1
	par := (i - 1) / 2
	for i > 0 && (*h)[i].dist > (*h)[par].dist {
		(*h)[i], (*h)[par] = (*h)[par], (*h)[i]
		i = par
		par = (i - 1) / 2
	}
}

func (h *pairHeap) pop() pair {
	if h.len() == 0 {
		return pair{}
	}
	p := (*h)[0]         // pop root
	size := h.len() - 1  // new size
	(*h)[0] = (*h)[size] // move last to root
	*h = (*h)[:size]     // shrink slice
	// percolate down while smaller than children
	i := 0
	for 2*i+1 < size { // while there is at least one child (left)
		par, lef, rig := i, 2*i+1, 2*i+2
		if lef < size && (*h)[par].dist < (*h)[lef].dist {
			par = lef
		}
		if rig < size && (*h)[par].dist < (*h)[rig].dist {
			par = rig
		}
		if par == i { // no swap happened
			break
		}
		(*h)[i], (*h)[par] = (*h)[par], (*h)[i] // now swap
		i = par
	}
	return p
}
