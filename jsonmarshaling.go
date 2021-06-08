package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Name string
	Body string `json:"body,omitempty"`
	Time int64
}

type Person struct {
	Name    string
	Age     int64
	Parents Parent
}

type Parent struct {
	Name     string
	Lastname string
}

func main() {
	var m = Message{"Hi", "", 0}
	byts, _ := json.Marshal(m)
	fmt.Println(string(byts))
	fmt.Println(m)
	var m2 Message

	var byts2 = []byte(`{"Name":"Hi","body":"Hello there","Time":234}`)
	json.Unmarshal(byts2, &m2)
	fmt.Println(m2)
	fmt.Println(m2.Body)
	fmt.Println(m2.Time)

	json1 := []byte(`{"Name":"Wednesday","Age":6,"Parents":{"Name":"Gomez","Lastname": "Morticia"}}`)
	var f interface{}
	json.Unmarshal(json1, &f)
	mymap := f.(map[string]interface{})

	//fmt.Println(mymap)
        var ap Person
	var parent Parent
	for k, v := range mymap {

	//A type switch is a construct that permits several type assertions in series.
	//A type switch is like a regular switch statement, but the cases in a type switch specify types (not values).. 
	//And those values are compared against the type of the value held by the given interface value.
		switch vv := v.(type) {
		case float64:
			fmt.Println(k, "is a float64 with value of ", vv)
		case string:
			fmt.Println(k, "is a string with value of ", vv)
		case []interface{}:
			for i, j := range vv {
				fmt.Println(i, "is member of array with value of", j)
			}
		//Incase type is a map[string]interface {}  unmarshall it to Parent struct
		case map[string]interface{}:
			jstr, _ := json.Marshal(v)
			json.Unmarshal(jstr, &parent)
			fmt.Println("Parent:" ,parent)	
		default:
			fmt.Println(k, "not sure what was it:", v)
		}

	}

	
	json.Unmarshal(json1, &ap)
	fmt.Println(ap)
	fmt.Println(ap.Parents.Name)
	fmt.Println(ap.Parents.Lastname)

}
