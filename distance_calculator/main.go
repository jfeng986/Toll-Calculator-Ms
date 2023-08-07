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

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewHTTPClient(aggEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
