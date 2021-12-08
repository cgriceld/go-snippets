package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func (out OutJ) Len() int {
	return len(out.Invalid)
}

func (out OutJ) Less(i, j int) bool {
	return out.CreatedInv[i].Before(out.CreatedInv[j])
}

func (out OutJ) Swap(i, j int) {
	out.CreatedInv[i], out.CreatedInv[j] = out.CreatedInv[j], out.CreatedInv[i]
	out.Invalid[i], out.Invalid[j] = out.Invalid[j], out.Invalid[i]
}

func getdata(path string) RawData {
	var d *json.Decoder
	if path != "" {
		raw, rerr := ioutil.ReadFile(path)
		if rerr != nil {
			log.Fatal(rerr)
		}
		d = json.NewDecoder(bytes.NewBuffer(raw))
	} else {
		d = json.NewDecoder(os.Stdin)
	}

	d.UseNumber()
	var billing RawData
	if derr := d.Decode(&billing); derr != nil {
		log.Fatal(derr)
	}

	return (billing)
}

func putdata(data Data) {
	fd, oerr := os.Create("out.json")
	if oerr != nil {
		log.Fatal(oerr)
	}
	defer func() {
		if cerr := fd.Close(); cerr != nil {
			log.Fatal(cerr)
		}
	}()

	fin := make([]OutJ, 0, len(data))
	for _, val := range data {
		sort.Sort(val)
		fin = append(fin, *val)
	}
	sort.Slice(fin, func(i, j int) bool { return fin[i].Company < fin[j].Company })

	out, merr := json.MarshalIndent(fin, "", "\t")
	if merr != nil {
		return
	}

	fmt.Fprintln(fd, string(out))
}
