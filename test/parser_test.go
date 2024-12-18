package test

import (
	"math"
	"os"
	"testing"

	"github.com/printchard/go-json-parser/parser"
)

const testCasesPath = "cases/"
const validCasesPath = testCasesPath + "valid/"
const invalidCasesPath = testCasesPath + "invalid/"

type testCase struct {
	path string
	err  bool
	val  any
}

var tests = []testCase{
	{
		path: validCasesPath + "arrays.json",
		err:  false,
		val:  map[string]any{"items": []any{"apple", "banana", "cherry"}},
	},
	{
		path: validCasesPath + "empty.json",
		err:  false,
		val:  map[string]any{},
	},
	{
		path: validCasesPath + "escaped.json",
		err:  false,
		val:  map[string]any{"quote": "\"Hello, world!\"", "path": "C:\\Program Files\\App"},
	},
	{
		path: validCasesPath + "mixed.json",
		err:  false,
		val: map[string]any{
			"numbers": []any{1, 2.5, -3},
			"flags":   []any{true, false, nil},
		},
	},
	{
		path: validCasesPath + "nested_objects.json",
		err:  false,
		val: map[string]any{
			"user": map[string]any{
				"name": "Bob",
				"address": map[string]any{
					"city": "Wonderland",
					"zip":  "12345",
				},
			},
		},
	},
	{
		path: validCasesPath + "simple_kv.json",
		err:  false,
		val:  map[string]any{"name": "Alice", "age": 25.0, "isStudent": false},
	},
	{
		path: invalidCasesPath + "duplicate_keys.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "invalid_number.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "invalid_root.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "missing_brace.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "missing_comma.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "nonstring_key.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "trailing_comma_array.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "trailing_comma.json",
		err:  true,
		val:  nil,
	},
	{
		path: invalidCasesPath + "unclosed_string.json",
		err:  true,
		val:  nil,
	},
}

func TestJson(t *testing.T) {
	for _, test := range tests {
		file, err := os.ReadFile(test.path)
		if err != nil {
			t.Errorf("Error reading test file %s: %s", test.path, err)
		}

		p := parser.New(string(file))
		parsedValue, err := p.Parse()
		if test.err {
			if err == nil {
				t.Errorf("Expected error but got <nil>: %v", parsedValue)
			}
			continue
		}

		switch val := test.val.(type) {
		case map[string]any:
			castValue, ok := parsedValue.(map[string]any)
			if !ok || !compareObj(val, castValue) {
				t.Errorf("Mismatched object, expected %v, got %v:", test.val, parsedValue)
			}
		case []any:
			castValue, ok := parsedValue.([]any)
			if !ok || !compareArr(val, castValue) {
				t.Errorf("Mismatched array, expected %v, got %v", test.val, parsedValue)
			}
		default:
			t.Errorf("Unexpected value: %v", parsedValue)
		}
	}
}

func compareObj(a, b map[string]any) bool {
	for k, v := range a {
		switch castV := v.(type) {
		case map[string]any:
			castB, ok := b[k].(map[string]any)
			if !ok || !compareObj(castV, castB) {
				return false
			}
		case []any:
			castB, ok := b[k].([]any)
			if !ok || !compareArr(castV, castB) {
				return false
			}
		case float64:
			castB, ok := b[k].(float64)
			if !ok || math.Abs(castV-castB) > 0.1 {
				return false
			}
		default:
			if v != b[k] {
				return false
			}
		}
	}
	return true
}

func compareArr(a, b []any) bool {
	for i, val := range a {
		switch castV := val.(type) {
		case map[string]any:
			castB, ok := b[i].(map[string]any)
			if !ok || !compareObj(castV, castB) {
				return false
			}
		case []any:
			castB, ok := b[i].([]any)
			if !ok || !compareArr(castV, castB) {
				return false
			}
		case float64:
		case int:
			castB, ok := b[i].(float64)
			if !ok || math.Abs(float64(castV)-castB) > 0.1 {
				return false
			}
		default:
			if val != b[i] {
				return false
			}
		}
	}
	return true
}
