package main

import (
	"bufio"
	"fmt"
	"os"

	"toy_lang/evaluator"
	"toy_lang/lexer"
	"toy_lang/parser"
)

func main() {

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
