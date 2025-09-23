package vm

import (
	"reflect"
	"testing"
	"toy_lang/compiler"
	"toy_lang/lexer"
	"toy_lang/parser"
)

type ramSlot struct {
	Addr  int
	Value any
}

type vmTest struct {
	input  string
	output []ramSlot
	vars   map[string]int
	id     int
}

func compareRamSlots(t *testing.T, wantI vmTest, got [1 << 11]any) {
	want := wantI.output

	for _, val := range want {
		gotVal := got[val.Addr]

		matched := false
		switch w := val.Value.(type) {
		case int:
			if g, ok := gotVal.(int); ok && g == w {
				matched = true
			}
		case float64:
			if g, ok := gotVal.(float64); ok && g == w {
				matched = true
			}
		case string:
			if g, ok := gotVal.(string); ok && g == w {
				matched = true
			}
		case bool:
			if g, ok := gotVal.(bool); ok && g == w {
				matched = true
			}
		default:
			// fallback to reflect if you hit an unexpected type
			if reflect.DeepEqual(gotVal, w) {
				matched = true
			}
		}

		if !matched {
			t.Errorf(
				"[Fail] Test %v failed at addr %d, wanted %v, got %v",
				wantI.id, val.Addr, val.Value, gotVal,
			)
		}
	}
}

func compareVMap(t *testing.T, wantI vmTest, got map[string]any) {
	if len(wantI.vars) != len(got) {
		t.Errorf("[FAIL] Not the same number of variables, wanted %d, got %d", len(wantI.vars), len(got))
	}
	for key, value := range wantI.vars {
		if got[key] != value {
			t.Errorf("[FAIL] Expected %v to be %v but it was %v", key, value, got[key])
		}
	}
}

func TestVM(t *testing.T) {
	tests := []vmTest{
		{
			input: "4",
			output: []ramSlot{
				{
					Addr:  0,
					Value: 4,
				},
			},
			id: 1,
		},
		{
			input: "4 + 3;",
			output: []ramSlot{
				{
					Addr:  2,
					Value: 7,
				},
			},
			id: 2,
		},
		{
			input: "4 + 2 - 2",
			output: []ramSlot{
				{
					Addr:  4,
					Value: 4,
				},
			},
			id: 3,
		},
		{
			input: "4 + 3 * 2;",
			output: []ramSlot{
				{
					Addr:  4,
					Value: 10,
				},
			},
			id: 4,
		},
		{
			input: "let x = 4;",
			vars: map[string]int{
				"x": 4,
			},
			id: 5,
		},
		{
			input: "let x = 5; x = 4;",
			vars: map[string]int{
				"x": 4,
			},
			id: 6,
		},
		{
			input: "let x = 4; let y = x * 3;",
			vars: map[string]int{
				"x": 4,
				"y": 12,
			},
			id: 7,
		},
	}
	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := parser.NewParser()
		compile := compiler.NewCompiler()
		vm := NewVm()

		gotR, gotV := vm.Execute(compile.Compile(parse.Parse(lex.Lex(tt.input))), false)
		if tt.output != nil {
			compareRamSlots(t, tt, *gotR)
		}
		if tt.vars != nil {
			compareVMap(t, tt, gotV)
		}
	}
}
