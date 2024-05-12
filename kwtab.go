package main

import (
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

type KwTabEntry[T any] struct {
	id  string
	sym T
}

type KwTab[T any] struct {
	entries []KwTabEntry[T]
	offsets []int
}

func NewKwTab[T any]() *KwTab[T] {
	return &KwTab[T]{}
}

func (t *KwTab[T]) Add(id string, v T) {
	t.entries = append(t.entries, KwTabEntry[T]{id, v})
}

func (t *KwTab[T]) Mark() {
	var dummy T
	t.offsets = append(t.offsets, len(t.entries))
	t.entries = append(t.entries, KwTabEntry[T]{"", dummy})
}

func (t *KwTab[T]) Lookup(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		k := t.offsets[n-1]
		for k < t.offsets[n] && id != t.entries[k].id {
			k++
		}
		if ok = k < t.offsets[n]; ok {
			v = t.entries[k].sym
		}
	}
	return
}

func (t *KwTab[T]) Lookup1(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		k := t.offsets[n-1]
		for ok = id != t.entries[k].id; ok; ok = id != t.entries[k].id {
			k++
			if k >= t.offsets[n] {
				return
			}
		}
		v = t.entries[k].sym
	}
	return
}

func (t *KwTab[T]) Lookup2(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		t.entries[t.offsets[n]].id = id
		k := t.offsets[n-1]
		for id != t.entries[k].id {
			k++
		}
		v = t.entries[k].sym
		ok = k < t.offsets[n]
	}

	return
}

func (t *KwTab[T]) Lookup3(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		k, u := t.offsets[n-1], t.offsets[n]
		for {
			if id == t.entries[k].id {
				break
			}
			if id == t.entries[k+1].id {
				k++
				break
			}
			k += 2
			if k >= u {
				return
			}
		}
		v = t.entries[k].sym
		ok = true
	}

	return
}

func (t *KwTab[T]) Lookup4(id string) (T, bool) {
	if n := len(id); 0 < n && n < 10 {
		s := t.entries[t.offsets[n-1]:t.offsets[n]]
		for _, e := range s {
			if e.id == id {
				return e.sym, true
			}
		}
	}
	var dummy T
	return dummy, false
}

type BinTab[T any] struct {
	bins [][]KwTabEntry[T]
}

func NewBinTab[T any](n int) *BinTab[T] {
	return &BinTab[T]{bins: make([][]KwTabEntry[T], n)}
}

func (t *BinTab[T]) Add(id string, v T) {
	t.bins[len(id)] = append(t.bins[len(id)], KwTabEntry[T]{id, v})
}

func (t *BinTab[T]) Mark() {
}

func (t *BinTab[T]) Lookup(id string) (T, bool) {
	if n := len(id); 0 < n && n < 10 {
		for _, e := range t.bins[n] {
			if e.id == id {
				return e.sym, true
			}
		}
	}
	var dummy T
	return dummy, false
}

type KvBinTab[T any] struct {
	kBins [][]string
	vBins [][]T
}

func NewKvBinTab[T any](n int) *KvBinTab[T] {
	return &KvBinTab[T]{kBins: make([][]string, n), vBins: make([][]T, n)}
}

func (t *KvBinTab[T]) Add(id string, v T) {
	t.kBins[len(id)] = append(t.kBins[len(id)], id)
	t.vBins[len(id)] = append(t.vBins[len(id)], v)
}

func (t *KvBinTab[T]) Mark() {
}

//go:noinline
func (t *KvBinTab[T]) Lookup(id string) (T, bool) {
	if n := len(id); 0 < n && n < 10 {
		for i, e := range t.kBins[n] {
			if e == id {
				return t.vBins[n][i], true
			}
		}
	}
	var dummy T
	return dummy, false
}

func initTab() *KwTab[Token] {
	t := NewKwTab[Token]()

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

func initKvBinTab() *KvBinTab[Token] {
	t := NewKvBinTab[Token](10)

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

var input []string

func init() {
	bs, _ := os.ReadFile("large.input")
	for _, s := range strings.Split(string(bs), "\n") {
		if len(s) > 0 {
			input = append(input, s)
		}
	}
}

func lex() int {
	// t := initTab()
	n := 0
	bt := initKvBinTab()
	for _, lex := range input {
		// fmt.Printf("%s: ", lex)
		if _, ok := bt.Lookup(lex); ok {
			n++
		}
		// fmt.Printf("%v %t\n", tok, ok)
	}
	return n
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println(lex())
}
