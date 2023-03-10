package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Had to "overload" the function to actually get what I wanted...
func getEnv(params ...string) string {
	value := os.Getenv(params[0])
	if len(value) == 0 {
		if len(params[1]) == 0 {
			panic("Could not find variable: " + params[0])
		}
		return params[1]
	}
	return value
}

func main() {

	// iNIT iNFLUX INFOS
	INFLUX_ADDRESS := getEnv("INFLUX_ADDRESS", "192.168.10.200")
	INFLUX_TOKEN := getEnv("INFLUX_TOKEN")
	INFLUX_PORT := getEnv("INFLUX_PORT", "8086")
	INFLUX_BUCKET_NAME := getEnv("INFLUX_BUCKET", "swe-ele-test")

	// INIT ELECTRICITY API INFOS
	API_URI := getEnv("API_URI", "https://mgrey.se/espot?format=json")

	data := getElectricityCost(API_URI)
	fmt.Printf("%+v\n", &data)

	// Create a new client using an InfluxDB server base URL and an authentication token
	c := influxdb2.NewClientWithOptions("http://"+INFLUX_ADDRESS+":"+INFLUX_PORT, INFLUX_TOKEN,
		influxdb2.DefaultOptions().SetBatchSize(20))

	writeDatabase(c, INFLUX_BUCKET_NAME, data)
	// Ensures background processes finishes
	defer c.Close()

}

func getElectricityCost(uri string) Data {
	spaceClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "github.com/boveloco/electricity-cost-sweden")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	var d Data

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &d)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return d
}

func writeDatabase(cli influxdb2.Client, bucket string, d Data) {

	writeAPI := cli.WriteAPIBlocking("bova", bucket)

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("region", "se1").
		AddField("price_sek", d.Se1[0].Price_sek).
		AddField("price_eur", d.Se1[0].Price_eur).
		AddField("kmeans", d.Se1[0].Kmeans).
		SetTime(time.Now())
	q := influxdb2.NewPointWithMeasurement("stat").
		AddTag("region", "se2").
		AddField("price_sek", d.Se2[0].Price_sek).
		AddField("price_eur", d.Se2[0].Price_eur).
		AddField("kmeans", d.Se2[0].Kmeans).
		SetTime(time.Now())

	r := influxdb2.NewPointWithMeasurement("stat").
		AddTag("region", "se3").
		AddField("price_sek", d.Se3[0].Price_sek).
		AddField("price_eur", d.Se3[0].Price_eur).
		AddField("kmeans", d.Se3[0].Kmeans).
		SetTime(time.Now())
	s := influxdb2.NewPointWithMeasurement("stat").
		AddTag("region", "se4").
		AddField("price_sek", d.Se4[0].Price_sek).
		AddField("price_eur", d.Se4[0].Price_eur).
		AddField("kmeans", d.Se4[0].Kmeans).
		SetTime(time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		panic(err)
	}
	err = writeAPI.WritePoint(context.Background(), q)
	if err != nil {
		panic(err)
	}
	err = writeAPI.WritePoint(context.Background(), r)
	if err != nil {
		panic(err)
	}
	err = writeAPI.WritePoint(context.Background(), s)
	if err != nil {
		panic(err)
	}

}
