package main

import (
	"log"

	"Toll-Calculator/aggregator/client"
)

const (
	kafkaTopic  = "obu-data"
	aggEndpoint = "http://localhost:30000/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	httpClient := client.NewHTTPClient(aggEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
