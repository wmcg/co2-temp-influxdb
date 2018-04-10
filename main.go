package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"

	"meter"

	"gopkg.in/alecthomas/kingpin.v2"

)

const (
	MyDB = "office_environment"
	username = ""
	password = ""
	)

var (
	device     = kingpin.Arg("device", "CO2 Meter device, such as /dev/hidraw2").Required().String()
	listenAddr = kingpin.Arg("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
)

func main() {

	//send_point()
	kingpin.Parse()
//	fmt.Printf("%v\n", measure())
	for {
		send_point(measure())
		time.Sleep(60000 * time.Millisecond)
	}
}

func send_point(co2_level int) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:80",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"location": "lon-office"}
	fields := map[string]interface{}{
		"ppm":  co2_level,
	}

	pt, err := client.NewPoint("co2_level", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func measure() (int) {
	meter := new(meter.Meter)
	err := meter.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
	}

	for {
		result, err := meter.Read()
		if err != nil {
			log.Fatalf("Something went wrong: '%v'", err)
		}
		// temperature.Set(result.Temperature)
		return result.Co2
	}
}