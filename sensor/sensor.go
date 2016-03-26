package main

import (
	"flag"
	"time"
	"strconv"
	"math/rand"
	"log"
)

var name = flag.String("name", "sensro", "Name of th sensor")
var freq = flag.Uint("freq", 5, "Updae frecuency in cycles/sec")
var max = flag.Float64("max", 5., "Maximum value for generated readings")
var min = flag.Float64("min", 1., "Minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "Maximum allowable change per measurement")

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = rand.Float64() * (*max-*min) + *min

var nominalValue = (*max - *min) / 2 + *min

func main() {
	flag.Parse()

	duration, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(duration)

	for range signal {
		calValue()
		log.Printf("Reading sent. Value: %v\n", value)

	}
}

func calValue() {
	var maxStep, minStep float64

	if value < nominalValue {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nominalValue - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nominalValue)
		minStep = -1 * *stepSize
	}

	value += rand.Float64() * (maxStep - minStep) + minStep
}