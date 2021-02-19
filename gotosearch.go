//Program sorts CS courses starting with ones that do not have any prereqs
package main

import (
	"fmt"
	"sort"
)
// Stores all the courses with prereqs of these course as a map
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":  {"discrete math"},
	"databases":        {"data structures"},
	"discrete math":    {"intro to programming"},
	"formal languages": {"discrete math"},
	"networks":         {"operating systems"},
	"operating systems": {"data structures",
		"computer organization"},
	"programming languages": {"data structures",
		"computer organization"},
}

func main() {
	// Here we declare VisitAll to later use it as recursive function within main
	var VisitAll func(v []string)
	// Order holds all the courses in order of prereq depency 
	var order []string
	
	// This map holds in memory all the courses we visited and later added to order
	var seen = make(map[string]bool)
	
	// Function visits all the siblings that have the prereq parent
	VisitAll = func(v []string) {
		
		// Sort all the siblings so that we have alphabetical order
		sort.Strings(v)
		for _, elem := range v {
			// If we haven't put the course in the list first use recursion to check if some courses 
			if !seen[elem] {
				seen[elem] = true
				VisitAll(prereqs[elem])
				// Once we fell through all the prereqs and came back add the course to the ordered slice 
				order = append(order, elem)
			}

		}

	}
	// Append the keys into the slice of strings 
	var keys []string
	for i := range prereqs {
		keys = append(keys, i)
	}

	VisitAll(keys)
	
	// Print out the results
	for i,j := range order{
		fmt.Println(i, " ", j)
	}
}
