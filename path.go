// Copyright 2012 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The bfspath package implements a breadth-first-search-based pathfinding 
// algorithm which runs faster than Dijkstra's algorithm in certain cases.
// This algorithm assumes that all edges have integral length. (Note that 
// this should never provide accuracy issues. For a graph in which distances 
// are stored as uint64, it is possible to measure edge-lengths up to 145
// trillion miles while maintaining accuracy of one inch.)
package bfspath

type Node interface {
	// Returns the number of edges
	// pointing away from this node
	Edges() int

	// Returns the node at the
	// end of the edge at index i
	Next(i int) Node

	// If the length of the edge
	// at index i is greater than 1,
	// decrement it by one and return
	// false. Else, return true
	Dec(i int) bool

	// Return whether or not this
	// node's previous node has
	// been decided yet
	PrevCheck() bool

	// Set this node's previous
	// node to n
	PrevSet(n Node)

	// Return this node's
	// previous node
	Prev() Node

	// Returns whether or not
	// this node is the
	// destination node
	IsEnd() bool

	// Returns whether or not
	// this node is the
	// start node
	IsStart() bool
}

/*
	The Prev and IsStart functions are 
	nowhere used in this code. However, 
	most reasonable uses of this package 
	require such functionality on the 
	part of the user, and so they are 
	included to save the user the headache 
	of type assertions.
*/

type llNode struct {
	Node
	next *llNode
}

type ll struct {
	head *llNode
	tail *llNode
}

func (l *ll) add(n Node) {
	newNode := new(llNode)
	newNode.Node = n

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
}

/*
	As specified in the documentation below,
	calling FindPath on a graph in which the
	start and end nodes are the same will cause
	an infinite loop. This problem is easily
	avoidable using a simple call to IsEnd().
	However, allowing for such behavior would
	put the burden on the user to make sure that
	when traversing their graph, they don't
	throw THEMSELVES into an infinite loop.
	Thus, we omit this behavior.
*/

// Takes start and end nodes, and returns
// path length; returns -1 if no path was found.
// Note that begining and ending nodes
// being the same will cause an infinite loop.
// 
// It is an invariant that if every node in
// the graph is of the same type when FindPath
// is called, no new types will be added.
// Most importantly, this means that any 
// type assertions which would have been 
// guaranted to succeed before FindPath 
// retain the guarantee after FindPath is called.
func FindPath(start Node) int {
	// Set start's previous pointer
	// because nil is a sentinal value.
	// This is never actually used, so 
	// it doesn't cause recursion problems
	start.PrevSet(start)

	listA := new(ll)
	listB := new(ll)
	listA.add(start)

	// Each iteration of BFS increases
	// path length by 1
	i := 1
END:
	for ; ; i++ {
		if listA.head == nil {
			return -1
		}
		for cur := listA.head; cur != nil; cur = cur.next {
			// Iterate through this node's adjacent nodes
			edges := cur.Node.Edges()
			added := false
			for j := 0; j < edges; j++ {
				next := cur.Node.Next(j)
				if cur.Node.Dec(j) && !next.PrevCheck() {
					next.PrevSet(cur.Node)
					if next.IsEnd() {
						break END
					}
					listB.add(next)
				} else if !added {
					listB.add(cur.Node)
					added = true
				}
			}
		}
		listA = listB
		listB = new(ll)
	}

	// Undo the change from
	// the beginning
	start.PrevSet(nil)
	return i
}
