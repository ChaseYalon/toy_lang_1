package main

import (
	"bufio"
	"fmt"
	"os"

	"toy_lang/compiler"
	"toy_lang/evaluator"
	"toy_lang/lexer"
	"toy_lang/parser"
	"toy_lang/vm"
)

func main() {
	if len(os.Args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print(">")
			scanner.Scan()
			text := scanner.Text()

			lex := lexer.NewLexer()
			toks := lex.Lex(text)
			parse := parser.NewParser()
			program := parse.Parse(toks)
			eval := evaluator.NewInterpreter()
			eval.Execute(program, true)
		}

	}
	if os.Args[1] == "--help" {
		fmt.Printf("Please call with the path to a .toy file or use with no path for a repl\n")
		os.Exit(0)
	}
	if os.Args[1] == "-c" {
		scanner := bufio.NewScanner(os.Stdin)

		for {

			fmt.Print(">")
			scanner.Scan()
			text := scanner.Text()
			lex := lexer.NewLexer()
			toks := lex.Lex(text)
			parse := parser.NewParser()
			program := parse.Parse(toks)
			compile := compiler.NewCompiler()
			bytecode := compile.Compile(program)
			vm := vm.NewVm()
			vm.Execute(bytecode, true)
		}
	}
	filePath := os.Args[1]

	source, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer()
	parse := parser.NewParser()
	program := parse.Parse(lex.Lex(string(source)))

	in := evaluator.NewInterpreter()

	in.Execute(program, false)
}
