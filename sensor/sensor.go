package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"time"
	"strconv"

	"github.com/AlejoJamC/distributedApp/dto"
	"github.com/AlejoJamC/distributedApp/qutils"
	"github.com/streadway/amqp"
)

// Location of the message broker's input listener
// TODO: Move this location to an external file of configuration
var url = "amqp://guest:guest@192.168.0.27/uri"


var freq = 		flag.Uint("freq", 5, "Updae frecuency in cycles/sec")
var max = 		flag.Float64("max", 5., "Maximum value for generated readings")
var min = 		flag.Float64("min", 1., "Minimum value for generated readings")
var name = 		flag.String("name", "sensro", "Name of th sensor")
var stepSize = 	flag.Float64("step", 0.1, "Maximum allowable change per measurement")


var nom = (*max - *min) / 2 + *min // Nominal values
var r = rand.New(rand.NewSource(time.Now().UnixNano())) // Random number
var value = r.Float64() * (*max-*min) + *min


func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer  conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)

	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"", // exchange
		sensorQueue.Name, // key
		false, // mandatory
		false, // immediate
		msg) // msg amqp.Publishing

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms") // duration

	signal := time.Tick(dur)

	buf := new(bytes.Buffer) // Buffer
	enc := gob.NewEncoder(buf) // Encoder

	for range signal {
		calValue()
		reading := dto.SensorMessage{
			Name: 		*name,
			Value: 		value,
			Timestamp: 	time.Now(),
		}

		buf.Reset() // Reset to initial position
		enc.Encode(reading) // Execute message encoding

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"", // Exchange name
			dataQueue.Name, // Key identificator of queue
			false, // mandatory
			false, // immediate
			msg) // msg amqp.Publishing

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