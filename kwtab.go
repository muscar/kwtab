package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

type Token int

const (
	TokIdent Token = iota
	TokKwIf
	TokKwDo
	TokKwOf
	TokKwOr
	TokKwTo
	TokKwIn
	TokKwIs
	TokKwBy
	TokKwEnd
	TokKwNil
	TokKwVar
	TokKwDiv
	TokKwMod
	TokKwFor
	TokKwElse
	TokKwThen
	TokKwTrue
	TokKwType
	TokKwCase
	TokKwElsif
	TokKwFalse
	TokKwArray
	TokKwBegin
	TokKwConst
	TokKwUntil
	TokKwWhile
	TokKwRecord
	TokKwRepeat
	TokKwReturn
	TokKwImport
	TokKwModule
	TokKwPointer
	TokKwProcedure
)

type BinTabEntry[T any] struct {
	key   string
	value T
}

type Bin[T any] []BinTabEntry[T]

type BinTab[T any] struct {
	bins []Bin[T]
}

func NewBinTab[T any](n int) *BinTab[T] {
	return &BinTab[T]{bins: make([]Bin[T], n)}
}

func (t *BinTab[T]) Add(key string, v T) {
	t.bins[len(key)] = append(t.bins[len(key)], BinTabEntry[T]{key, v})
}

func (t *BinTab[T]) Mark() {
}

func (t *BinTab[T]) Lookup(key string) (v T, ok bool) {
	if n := len(key); n < len(t.bins) {
		for _, e := range t.bins[n] {
			if e.key == key {
				return e.value, true
			}
		}
	}
	return
}

func initBinTab() *BinTab[Token] {
	t := NewBinTab[Token](10)

	t.Mark()
	t.Mark()
	t.Add("IF", TokKwIf)
	t.Add("DO", TokKwDo)
	t.Add("OF", TokKwOf)
	t.Add("OR", TokKwOr)
	t.Add("TO", TokKwTo)
	t.Add("IN", TokKwIn)
	t.Add("IS", TokKwIs)
	t.Add("BY", TokKwBy)
	t.Mark()
	t.Add("END", TokKwEnd)
	t.Add("NIL", TokKwNil)
	t.Add("VAR", TokKwVar)
	t.Add("DIV", TokKwDiv)
	t.Add("MOD", TokKwMod)
	t.Add("FOR", TokKwFor)
	t.Mark()
	t.Add("ELSE", TokKwElse)
	t.Add("THEN", TokKwThen)
	t.Add("TRUE", TokKwTrue)
	t.Add("TYPE", TokKwType)
	t.Add("CASE", TokKwCase)
	t.Mark()
	t.Add("ELSIF", TokKwElsif)
	t.Add("FALSE", TokKwFalse)
	t.Add("ARRAY", TokKwArray)
	t.Add("BEGIN", TokKwBegin)
	t.Add("CONST", TokKwConst)
	t.Add("UNTIL", TokKwUntil)
	t.Add("WHILE", TokKwWhile)
	t.Mark()
	t.Add("RECORD", TokKwRecord)
	t.Add("REPEAT", TokKwRepeat)
	t.Add("RETURN", TokKwReturn)
	t.Add("IMPORT", TokKwImport)
	t.Add("MODULE", TokKwModule)
	t.Mark()
	t.Add("POINTER", TokKwPointer)
	t.Mark()
	t.Mark()
	t.Add("PROCEDURE", TokKwProcedure)
	t.Mark()

	return t
}

var kwMap = map[string]Token{
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

var bytesInput [][]byte

var input []string

func initInputs(path string) {
	bs, _ := os.ReadFile(path)
	for _, s := range bytes.Split(bs, []byte{'\n'}) {
		if len(s) > 0 {
			bytesInput = append(bytesInput, s)
		}
	}

	for _, s := range strings.Split(string(bs), "\n") {
		if len(s) > 0 {
			input = append(input, s)
		}
	}
}

func lex[T any](bt *BinTab[T]) int {
	n := 0
	for _, lex := range input {
		if _, ok := bt.Lookup(lex); ok {
			n++
		}
	}
	return n
}

func lexM[T any](m map[string]T) int {
	n := 0
	for _, lex := range input {
		if _, ok := m[lex]; ok {
			n++
		}
	}
	return n
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	initInputs("larger.input")
	bt := initBinTab()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println(lex(bt))
	fmt.Println(lexM(kwMap))
}
