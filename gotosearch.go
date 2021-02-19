package main

import (
	"fmt"
	"sort"
)

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
	var VisitAll func(v []string)
	
	var seen = make(map[string]bool)
	var order []string

	VisitAll = func(v []string) {

		sort.Strings(v)
		for _, elem := range v {
			if !seen[elem] {
				seen[elem] = true
				VisitAll(prereqs[elem])
				order = append(order, elem)
			}

		}

	}
	var keys []string

	for i := range prereqs {
		keys = append(keys, i)

	}

	VisitAll(keys)
	for _,i := range order{
	fmt.Println(i)
	}

}
