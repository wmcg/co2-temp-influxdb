package main

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"

	"meter"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	MyDB     = "office_environment"
	username = ""
	password = ""
)

var (
	device    = kingpin.Arg("device", "CO2 Meter device, such as /dev/hidraw2").Required().String()
	influxUrl = kingpin.Arg("influx-url", "The address and port of the influx server - localhost:8088").String()
	influxDb  = kingpin.Arg("influx-db", "The influxdb database. e.g. primary_db").Default("http://localhost:8088").String()
)

func main() {

	//send_point()
	kingpin.Parse()
	//	fmt.Printf("%v\n", measure())
	client := create_influx_client(influxUrl, influxDb, "", "")
	for {
		co2, temp := measure()
		send_points(client, co2, temp)
		time.Sleep(60000 * time.Millisecond)
	}
}

func create_influx_client(url *string, db *string, username string, password string) client.Client {
	// Create a new influx client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + *url,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func send_points(c client.Client, co2_level int, temp float64) {

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  *influxDb,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"location": "lon-office"}
	fields := map[string]interface{}{
		"ppm":  co2_level,
		"temp": temp,
	}

	pt1, err := client.NewPoint("levels", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt1)

	// Write the batch
	if err := c.Write(bp); err != nil {
		println("Couldnt write point")
		// log.Warn(err)
	}
}

func measure() (int, float64) {
	meter := new(meter.Meter)
	err := meter.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
	}

	for {
		result, err := meter.Read()
		if err != nil {
			log.Fatalf("Error in Reading from meter: '%v'", err)
		}
		// temperature.Set(result.Temperature)
		return result.Co2, result.Temperature
	}
}
