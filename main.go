package main

import (
	"fmt"
	"os"

	"github.com/printchard/go-json-parser/parser"
)

func main() {
	json := `{"name": "Alice", "age": 25, "isStudent": false, "hobbies": ["reading", "hiking"]}`

	parser := parser.New(json)
	tokens, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error parsing the json: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(tokens)
}
