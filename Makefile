obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

rec:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

cal:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

agg:
	@go build -o bin/agg ./aggregator
	@./bin/agg

.PHONY: obu