package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/ssimunic/gosensors"
)

func main() {

	influxUrl := flag.String("influxUrl", "", "")
	influxAuth := flag.String("influxAuth", "", "")
	influxTable := flag.String("influxTable", "", "")
	computerName := flag.String("computerName", "", "")

	if len(os.Args) < 4 {
		fmt.Println("expected 'influxUrl', 'influxAuth', 'influxTable', 'computerName' arguments")
		os.Exit(1)
	}

	s, err := gosensors.NewFromSystem()
	if err != nil {
		panic(err)
	}

	client := influxdb2.NewClient(influxUrl, influxAuth)
	writeAPI := client.WriteAPIBlocking("", influxTable)

	var CPUTemp, SSDTemp float64

	for chip := range s.Chips {
		if chip == "k10temp-pci-00c3" {
			for key, value := range s.Chips[chip] {
				if key == "temp1" {
					cel := substr(value, 1, 4)
					CPUTemp = toFahrenheit(cel)
					fmt.Println(CPUTemp)
				}
			}
		}
		if chip == "nvme-pci-0100" {
			for key, value := range s.Chips[chip] {
				if key == "Composite" {
					cel := substr(value, 1, 4)
					SSDTemp = toFahrenheit(cel)
					fmt.Println(SSDTemp)
				}
			}
		}
	}

	recordedDate := time.Now().UTC()
	p := influxdb2.NewPointWithMeasurement("lenovo").
		AddField("cpu_temperature", CPUTemp).
		AddField("ssd_temperature", SSDTemp).
		AddField("computer_name", computerName).
		SetTime(recordedDate.Date)

	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Printf("Write error: %s\n", err.Error())
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

func toFahrenheit(input string) float64 {
	f, err := strconv.ParseFloat(input, 8)
	if err != nil {
		fmt.Println(err)
	}
	conv := (f * 1.8000) + 32
	return conv
}
