package main

import (
	"fmt"

	"github.com/printchard/go-json-parser/parser"
)

type Pet struct {
	Name    string `json:"name"`
	Species string `json:"species"`
}

type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Pet     Pet      `json:"pet"`
	Hobbies []string `json:"hobbies"`
}

func main() {
	myJson := `
		[
			{
				"name": "Alice",
				"age": 25,
				"isStudent": false,
				"pet": {
					"name": "Bob",
					"species": "dog"
				},
				"hobbies": ["reading", "hiking"]
			},
			{
				"name": "Charlie",
				"age": 30,
				"isStudent": true,
				"pet": {
					"name": "Daisy",
					"species": "cat"
				},
				"hobbies": ["gaming", "cooking"]
			}
		]
	`

	parser := parser.New(myJson)
	people := []Person{}
	if err := parser.ParseInto(&people); err != nil {
		panic(err)
	}

	for _, person := range people {
		fmt.Printf("%+v\n", person)
	}
}
