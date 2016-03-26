package main

import (
	"flag"
	"time"
	"strconv"
	"math/rand"
	"log"
	"encoding/gob"
	"github.com/AlejoJamC/distributedApp/dto"
	"bytes"
)

// Location of the message broker's input listener
// TODO: Move this location to an external file of configuration
var url = "amqp://guest:guest@localhost:5672"


var name = flag.String("name", "sensro", "Name of th sensor")
var freq = flag.Uint("freq", 5, "Updae frecuency in cycles/sec")
var max = flag.Float64("max", 5., "Maximum value for generated readings")
var min = flag.Float64("min", 1., "Minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "Maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano())) // Random number

var value = r.Float64() * (*max-*min) + *min

var nom = (*max - *min) / 2 + *min // Nominal values

func main() {
	flag.Parse()

	duration, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(duration)

	buf := new(bytes.Buffer) // Buffer
	enc := gob.NewEncoder(buf) // Encoder

	for range signal {
		calValue()
		reading := dto.SensorMessage{
			Name: *name,
			Value: value,
			Timestamp: time.Now(),
		}

		buf.Reset() // Reset to initial position
		enc.Encode(reading) // Execute message encoding

		log.Printf("Reading sent. Value: %v\n", value)

	}
}

func calValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += rand.Float64() * (maxStep - minStep) + minStep
}