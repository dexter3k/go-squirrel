package squirrel

import ()

type tableHashNode struct {
}

type Table struct {
	nodes []tableHashNode
}

func NewTable(initialSize uint) *Table {
	nodeCount := uint(4)
	for nodeCount < initialSize {
		nodeCount <<= 1
	}
	return &Table{
		nodes: make([]tableHashNode, nodeCount),
	}
}
