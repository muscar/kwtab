package main

import (
	"testing"
)

var inputPath = "larger.input"

func BenchmarkBinTab(b *testing.B) {
	initInputs(inputPath)
	t := initBinTab()
	b.ResetTimer()

	n := 0
	for i := 0; i < b.N; i++ {
		for _, lex := range input {
			if _, ok := t.Lookup(lex); ok {
				n++
			}
		}
	}
}

func BenchmarkMap(b *testing.B) {
	initInputs(inputPath)
	t := map[string]Token{
		"IF":        TokKwIf,
		"DO":        TokKwDo,
		"OF":        TokKwOf,
		"OR":        TokKwOr,
		"TO":        TokKwTo,
		"IN":        TokKwIn,
		"IS":        TokKwIs,
		"BY":        TokKwBy,
		"END":       TokKwEnd,
		"NIL":       TokKwNil,
		"VAR":       TokKwVar,
		"DIV":       TokKwDiv,
		"MOD":       TokKwMod,
		"FOR":       TokKwFor,
		"ELSE":      TokKwElse,
		"THEN":      TokKwThen,
		"TRUE":      TokKwTrue,
		"TYPE":      TokKwType,
		"CASE":      TokKwCase,
		"ELSIF":     TokKwElsif,
		"FALSE":     TokKwFalse,
		"ARRAY":     TokKwArray,
		"BEGIN":     TokKwBegin,
		"CONST":     TokKwConst,
		"UNTIL":     TokKwUntil,
		"WHILE":     TokKwWhile,
		"RECORD":    TokKwRecord,
		"REPEAT":    TokKwRepeat,
		"RETURN":    TokKwReturn,
		"IMPORT":    TokKwImport,
		"MODULE":    TokKwModule,
		"POINTER":   TokKwPointer,
		"PROCEDURE": TokKwProcedure,
	}
	b.ResetTimer()

	n := 0
	for i := 0; i < b.N; i++ {
		for _, lex := range input {
			if _, ok := t[lex]; ok {
				n++
			}
		}
	}
}
