package main

import (
	"fmt"
	"strings"
)

//define type of a function which takes string as an argument and returns string
type strfunc func(string) string

//user created function with type of 'strfunc'
var upper strfunc = strfunc(func(input string) string {
	return strings.ToUpper(input)
})

//function that will execute and print a list of functions with type of 'strfunc'
func myFunc(spec string, funcs []strfunc) {
	for _, exec := range funcs {
		fmt.Println(exec(spec))
	}
}

func main() {
	myFunc("Hola", []strfunc{upper, strings.ToUpper, strings.ToLower})
}
