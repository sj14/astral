// usage: citycsv2go.go <infile> <outfile>
// reads a .csv.gz file (as being downloadable from
// https://simplemaps.com/data/world-cities ) and
// generates a city.go file which provides a slice
// named `cities` afterwards

package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
)

type city struct {
	Lat  float64
	Long float64
	CC   [2]byte
	Name string
}

func main() {

	fin := os.Stdin
	if len(os.Args) > 1 {
		fin, _ = os.Open(os.Args[1])
	}
	fout := os.Stdout
	if len(os.Args) > 2 {
		fout, _ = os.Create(os.Args[2])
	}

	gzr, _ := gzip.NewReader(fin)
	r := csv.NewReader(gzr)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cities := []city{}
	for _, row := range records[1:] {
		lat, _ := strconv.ParseFloat(row[2], 64)
		long, _ := strconv.ParseFloat(row[3], 64)
		c := city{
			Lat:  lat,
			Long: long,
			Name: row[1],
		}
		c.CC[0], c.CC[1] = row[5][0], row[5][1]
		cities = append(cities, c)
	}

	renderCitiesGoCode2(fout, cities, "main")
}

func renderCitiesGoCode(cities []city, pkg string) {
    fmt.Println(`// note: this is an generated file, do not
// modify it.

package`, pkg, `

type city struct {
	Lat  float64
	Long float64
	CC   [2]byte
	Name string
}

var cities = []city{`)
	for _, c := range cities {
		fmt.Printf("\t{%f,%f,{%#v,%#v},%q},\n", c.Lat, c.Long, c.CC[0], c.CC[1], c.Name)
	}
	fmt.Println(`}`)
}

func renderCitiesGoCode2(w io.Writer, cities []city, pkg string) {
	fmt.Fprintln(w, `// note: this is an generated file, do not
// modify it.

package`, pkg, `

type city struct {
	Lat    float64
	Long   float64
	NameCC string
}

var cities = []city{`)
	for _, c := range cities {
		fmt.Fprintf(w, "\t{%f,%f,\"%s,%c%c\"},\n", c.Lat, c.Long, c.Name, c.CC[0], c.CC[1])
	}
	fmt.Fprintln(w, `}`)
}

func renderCitiesCompressedGOB(cities []city) {
	buf := bytes.NewBuffer(nil)
	flateWriter, _ := flate.NewWriter(buf, 9)
	enc := gob.NewEncoder(flateWriter)
	enc.Encode(cities)
	fmt.Printf("%#v", buf.Bytes())
}

