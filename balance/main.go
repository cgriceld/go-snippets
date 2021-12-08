package main

import (
	"flag"
	"os"
	"reflect"
	"time"
)

type OutJ struct {
	Company    string        `json:"company"`
	Valid      int           `json:"valid_operations_count"`
	Balance    int           `json:"balance"`
	Invalid    []interface{} `json:"invalid_operations,omitempty"`
	CreatedInv []time.Time   `json:"-"`
}

type Nested struct {
	Typ     interface{} `json:"type"`
	Value   interface{} `json:"value"`
	ID      interface{} `json:"id"`
	Created interface{} `json:"created_at"`
}

type Raw struct {
	Company interface{} `json:"company"`
	Oper    Nested      `json:"operation"`
	Nested
}

type Parse struct {
	raw    *Raw
	create time.Time
	id     interface{}
	typ    string
	value  int
}

type RawData []Raw

type Data map[string]*OutJ

type Checker func(interface{}, *Parse) bool

// searches for passed field at the root, then in nested field
func (r *Raw) getField(field string) interface{} {
	var f reflect.Value
	if f = reflect.ValueOf(r).Elem().FieldByName(field); f.IsNil() {
		if f = reflect.ValueOf(&r.Oper).Elem().FieldByName(field); f.IsNil() {
			return nil
		}
	}
	return f.Interface()
}

func main() {
	data := make(Data)

	fptr := flag.String("file", "", "path to input file")
	flag.Parse()

	if *fptr != "" {
		parse(data, getdata(*fptr))
	} else if path, ok := os.LookupEnv("FILE"); ok {
		parse(data, getdata(path))
	} else {
		parse(data, getdata(""))
	}

	putdata(data)
}
