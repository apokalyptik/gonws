package main

import (
	"flag"
	"log"

	"github.com/apokalyptik/gonws/forecast"
)

func main() {
	var lat float64
	var lon float64
	var metric bool
	var language string
	flag.Float64Var(&lat, "lat", 0, "Latitude")
	flag.Float64Var(&lon, "lon", 0, "Longitude")
	flag.BoolVar(&metric, "metric", false, "false: imperial, true: metric")
	flag.StringVar(&language, "language", "english", "written to request")
	flag.Parse()
	c := forecast.New().Lang(language)
	if metric {
		c.Metric()
	}
	r, err := c.Get(lat, lon)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", r)
}
