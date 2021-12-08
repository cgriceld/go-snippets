package main

import (
	"encoding/json"
	"math"
	"strconv"
	"time"
)

func checkID(val interface{}, tmp *Parse) bool {
	switch i := val.(type) {
	case json.Number:
		if j, err := i.Int64(); err == nil {
			tmp.id = j
			return true
		}
	case string:
		tmp.id = i
		return true
	}
	return false
}

func checkType(val interface{}, tmp *Parse) bool {
	var ok bool
	if tmp.typ, ok = val.(string); ok {
		switch tmp.typ {
		case
			"income",
			"outcome",
			"+",
			"-":
			return true
		default:
			return false
		}
	}
	return false
}

func checkValue(val interface{}, tmp *Parse) bool {
	switch v := val.(type) {
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			tmp.value = i
			return true
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			if f == math.Trunc(f) {
				tmp.value = int(f)
				return true
			}
		}

	case json.Number:
		if i, err := v.Int64(); err == nil {
			tmp.value = int(i)
			return true
		}
		if f, err := v.Float64(); err == nil {
			if f == math.Trunc(f) {
				tmp.value = int(f)
				return true
			}
		}
	}
	return false
}

func checkCreate(val interface{}, tmp *Parse) bool {
	if t, ok := val.(string); ok {
		var err error
		if tmp.create, err = time.Parse(time.RFC3339, t); err == nil {
			return true
		}
	}
	return false
}

// gets field with passed param and if it exists applies check function
func check(param string, f Checker, tmp *Parse, currComp *OutJ) bool {
	val := tmp.raw.getField(param)
	if val == nil || (val != nil && !f(val, tmp)) {
		if currComp != nil {
			currComp.Invalid = append(currComp.Invalid, tmp.id)
			currComp.CreatedInv = append(currComp.CreatedInv, tmp.create)
		}
		return true
	}
	return false
}

func parse(data Data, billing RawData) {
	var tmp Parse
	for _, curr := range billing {
		if curr.Company == nil {
			continue
		}
		currComp, ok := curr.Company.(string)
		if !ok {
			continue
		}

		tmp = Parse{raw: &curr}

		if check("Created", checkCreate, &tmp, nil) ||
			check("ID", checkID, &tmp, nil) {
			continue
		}

		if _, ok := data[currComp]; !ok {
			data[currComp] = &OutJ{Company: currComp}
		}

		if check("Typ", checkType, &tmp, data[currComp]) ||
			check("Value", checkValue, &tmp, data[currComp]) {
			continue
		}

		data[currComp].Valid++
		if tmp.typ == "income" || tmp.typ == "+" {
			data[currComp].Balance += tmp.value
		} else {
			data[currComp].Balance -= tmp.value
		}
	}
}
