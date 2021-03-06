// Copyright 2012 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/joshlf13/bfspath"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <filename>", os.Args[0])
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error openning file: %v\n", err)
		os.Exit(2)
	}
	start, end := importGraph(file)
	length := bfspath.FindPath(start)
	if length == -1 {
		fmt.Println("That graph has no path between the endpoints!")
		return
	}
	fmt.Printf("Shortest path has length %v\n", length)
	for cur := bfspath.Node(end); !cur.IsStart(); cur = cur.Prev() {
		c, _ := cur.(*node)
		fmt.Printf("%v <- ", c.name)
	}
	fmt.Printf("%v\n", start.name)
}
