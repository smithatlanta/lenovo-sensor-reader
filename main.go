package main

import (
	"fmt"
	"github.com/ssimunic/gosensors"
	"strconv"
)

func main() {
	s, err := gosensors.NewFromSystem()
	if err != nil {
		panic(err)
	}

	// Iterate over chips
	for chip := range s.Chips {
		if chip == "coretemp-isa-0000" {
			for key, value := range s.Chips[chip] {
				if key == "Package id 0" {
					cel := substr(value, 1, 4)
					fahr := toFahrenheit(cel)
					fmt.Println(fahr)
				}
			}
		} else if chip == "k10temp-pci-00c3" {
			for key, value := range s.Chips[chip] {
				if key == "temp1" {
					cel := substr(value, 1, 4)
					fahr := toFahrenheit(cel)
					fmt.Println(fahr)
				}
			}
		}
	}
}

func substr(input string, start int, length int) string {
    asRunes := []rune(input)

    if start >= len(asRunes) {
        return ""
    }

    if start+length > len(asRunes) {
        length = len(asRunes) - start
    }

    return string(asRunes[start : start+length])
}

func toFahrenheit(input string) float64{
	f, err := strconv.ParseFloat(input, 8)
	if err != nil {
		fmt.Println(err)
	} 
	conv := (f*1.8000) + 32
	return conv
}
