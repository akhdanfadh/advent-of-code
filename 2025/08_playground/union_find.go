package main

type circuits struct {
	parent []int
	size   []int
}

func initCircuits(size int) *circuits {
	parent := make([]int, size)
	sizes := make([]int, size)
	for i := range size {
		parent[i] = i
		sizes[i] = 1 // each start with its own group (size 1)
	}
	return &circuits{parent: parent, size: sizes}
}

func (c *circuits) getRootSizes() []int {
	seen := make(map[int]struct{})
	sizes := make([]int, 0)
	for i := range c.parent {
		root := c.find(i)
		if _, exists := seen[root]; !exists {
			seen[root] = struct{}{}
			sizes = append(sizes, c.size[root]) // get size from root
		}
	}
	return sizes
}

func (c *circuits) find(pid int) int {
	// keep following parent until root (point to itself)
	root := pid
	for root != c.parent[root] {
		root = c.parent[root]
	}
	// path compression (flatten tree)
	for pid != root {
		next := c.parent[pid]
		c.parent[pid] = root
		pid = next
	}
	return root
}

func (c *circuits) union(pid1, pid2 int) {
	// attach one root to another
	root1 := c.find(pid1)
	root2 := c.find(pid2)
	if root1 != root2 {
		// attach smaller to larger
		if c.size[root1] < c.size[root2] {
			c.parent[root1] = root2
			c.size[root2] += c.size[root1]
		} else {
			c.parent[root2] = root1
			c.size[root1] += c.size[root2]
		}
	}
}
