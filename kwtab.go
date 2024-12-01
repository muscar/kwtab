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

type KwTab[T any] struct {
	entries []BinTabEntry[T]
	offsets []int
}

func NewKwTab[T any]() *KwTab[T] {
	return &KwTab[T]{}
}

func (t *KwTab[T]) Add(id string, v T) {
	t.entries = append(t.entries, BinTabEntry[T]{id, v})
}

func (t *KwTab[T]) Mark() {
	var dummy T
	t.offsets = append(t.offsets, len(t.entries))
	t.entries = append(t.entries, BinTabEntry[T]{"", dummy})
}

func (t *KwTab[T]) Lookup(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < len(t.offsets) {
		k := t.offsets[n-1]
		for k < t.offsets[n] && id != t.entries[k].key {
			k++
		}
		if ok = k < t.offsets[n]; ok {
			v = t.entries[k].value
		}
	}
	return
}

func (t *KwTab[T]) Lookup1(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		k := t.offsets[n-1]
		for ok = id != t.entries[k].key; ok; ok = id != t.entries[k].key {
			k++
			if k >= t.offsets[n] {
				return
			}
		}
		v = t.entries[k].value
	}
	return
}

func (t *KwTab[T]) Lookup2(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		t.entries[t.offsets[n]].key = id
		k := t.offsets[n-1]
		for id != t.entries[k].key {
			k++
		}
		v = t.entries[k].value
		ok = k < t.offsets[n]
	}

	return
}

func (t *KwTab[T]) Lookup3(id string) (v T, ok bool) {
	if n := len(id); 0 < n && n < 10 {
		k, u := t.offsets[n-1], t.offsets[n]
		for {
			if id == t.entries[k].key {
				break
			}
			if id == t.entries[k+1].key {
				k++
				break
			}
			k += 2
			if k >= u {
				return
			}
		}
		v = t.entries[k].value
		ok = true
	}

	return
}

func (t *KwTab[T]) Lookup4(id string) (T, bool) {
	if n := len(id); 0 < n && n < 10 {
		s := t.entries[t.offsets[n-1]:t.offsets[n]]
		for _, e := range s {
			if e.key == id {
				return e.value, true
			}
		}
	}
	var dummy T
	return dummy, false
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

type KvBinTab[T any] struct {
	kBins [][][]byte
	vBins [][]T
}

func NewKvBinTab[T any](n int) *KvBinTab[T] {
	return &KvBinTab[T]{kBins: make([][][]byte, n), vBins: make([][]T, n)}
}

func (t *KvBinTab[T]) Add(id []byte, v T) {
	t.kBins[len(id)] = append(t.kBins[len(id)], id)
	t.vBins[len(id)] = append(t.vBins[len(id)], v)
}

func (t *KvBinTab[T]) Mark(n int) {
	if len(t.kBins[n]) == 0 || len(t.kBins[n])%2 == 1 {
		var dummy T
		t.kBins[n] = append(t.kBins[n], nil)
		t.vBins[n] = append(t.vBins[n], dummy)
	}
}

//go:noinline
func (t *KvBinTab[T]) Lookup(id []byte) (v T, ok bool) {
	if n := len(id); n < 10 {
		bin := t.kBins[n]
		bin[len(bin)-1] = id
		for i := 0; ; i += 2 {
			if bytes.Equal(bin[i], id) {
				return t.vBins[n][i], i < len(bin)-1
			}
			if bytes.Equal(bin[i+1], id) {
				i++
				return t.vBins[n][i], i < len(bin)-1
			}
		}
	}
	return
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

	t.Mark(0)
	t.Mark(1)
	t.Add([]byte("IF"), TokKwIf)
	t.Add([]byte("DO"), TokKwDo)
	t.Add([]byte("OF"), TokKwOf)
	t.Add([]byte("OR"), TokKwOr)
	t.Add([]byte("TO"), TokKwTo)
	t.Add([]byte("IN"), TokKwIn)
	t.Add([]byte("IS"), TokKwIs)
	t.Add([]byte("BY"), TokKwBy)
	t.Mark(2)
	t.Add([]byte("END"), TokKwEnd)
	t.Add([]byte("NIL"), TokKwNil)
	t.Add([]byte("VAR"), TokKwVar)
	t.Add([]byte("DIV"), TokKwDiv)
	t.Add([]byte("MOD"), TokKwMod)
	t.Add([]byte("FOR"), TokKwFor)
	t.Mark(3)
	t.Add([]byte("ELSE"), TokKwElse)
	t.Add([]byte("THEN"), TokKwThen)
	t.Add([]byte("TRUE"), TokKwTrue)
	t.Add([]byte("TYPE"), TokKwType)
	t.Add([]byte("CASE"), TokKwCase)
	t.Mark(4)
	t.Add([]byte("ELSIF"), TokKwElsif)
	t.Add([]byte("FALSE"), TokKwFalse)
	t.Add([]byte("ARRAY"), TokKwArray)
	t.Add([]byte("BEGIN"), TokKwBegin)
	t.Add([]byte("CONST"), TokKwConst)
	t.Add([]byte("UNTIL"), TokKwUntil)
	t.Add([]byte("WHILE"), TokKwWhile)
	t.Mark(5)
	t.Add([]byte("RECORD"), TokKwRecord)
	t.Add([]byte("REPEAT"), TokKwRepeat)
	t.Add([]byte("RETURN"), TokKwReturn)
	t.Add([]byte("IMPORT"), TokKwImport)
	t.Add([]byte("MODULE"), TokKwModule)
	t.Mark(6)
	t.Add([]byte("POINTER"), TokKwPointer)
	t.Mark(7)
	t.Mark(8)
	t.Add([]byte("PROCEDURE"), TokKwProcedure)
	t.Mark(9)

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
	// t := initTab()
	n := 0
	for _, lex := range input {
		// fmt.Printf("%s: ", lex)
		if _, ok := bt.Lookup(lex); ok {
			n++
		}
		// fmt.Printf("%v %t\n", tok, ok)
	}
	return n
}

func lexM[T any](m map[string]T) int {
	// t := initTab()
	n := 0
	for _, lex := range input {
		// fmt.Printf("%s: ", lex)
		if _, ok := m[lex]; ok {
			n++
		}
		// fmt.Printf("%v %t\n", tok, ok)
	}
	return n
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	initInputs("larger.input")
	bt := initBinTab()
	// m := map[string]Token{
	// 	"IF":        TokKwIf,
	// 	"DO":        TokKwDo,
	// 	"OF":        TokKwOf,
	// 	"OR":        TokKwOr,
	// 	"TO":        TokKwTo,
	// 	"IN":        TokKwIn,
	// 	"IS":        TokKwIs,
	// 	"BY":        TokKwBy,
	// 	"END":       TokKwEnd,
	// 	"NIL":       TokKwNil,
	// 	"VAR":       TokKwVar,
	// 	"DIV":       TokKwDiv,
	// 	"MOD":       TokKwMod,
	// 	"FOR":       TokKwFor,
	// 	"ELSE":      TokKwElse,
	// 	"THEN":      TokKwThen,
	// 	"TRUE":      TokKwTrue,
	// 	"TYPE":      TokKwType,
	// 	"CASE":      TokKwCase,
	// 	"ELSIF":     TokKwElsif,
	// 	"FALSE":     TokKwFalse,
	// 	"ARRAY":     TokKwArray,
	// 	"BEGIN":     TokKwBegin,
	// 	"CONST":     TokKwConst,
	// 	"UNTIL":     TokKwUntil,
	// 	"WHILE":     TokKwWhile,
	// 	"RECORD":    TokKwRecord,
	// 	"REPEAT":    TokKwRepeat,
	// 	"RETURN":    TokKwReturn,
	// 	"IMPORT":    TokKwImport,
	// 	"MODULE":    TokKwModule,
	// 	"POINTER":   TokKwPointer,
	// 	"PROCEDURE": TokKwProcedure,
	// }

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
	// fmt.Println(lexM(m))
}
