//go:generate go run ../../data/citycsv2go.go ../../data/worldcities.csv.gz cities_gen.go

package main

import (
	"fmt"
	"strings"
)

func latLongFromCity(city string) (lat, long float64, err error) {

	city = strings.TrimSpace(city)
	city = strings.ToLower(city)
	for _, c := range cities {
		ncc := strings.ToLower(c.NameCC)
		if strings.HasPrefix(ncc, city) {
			return c.Lat, c.Long, nil
		}
		// "new york" -> "newyork"
		ncc = strings.ReplaceAll(ncc, " ", "")
		if strings.HasPrefix(ncc, city) {
			return c.Lat, c.Long, nil
		}
	}

	return 0, 0, fmt.Errorf("not found")
}

func printCities() {
	for _, c := range cities {
		fmt.Println(c.NameCC, c.Lat, c.Long)
	}
}
