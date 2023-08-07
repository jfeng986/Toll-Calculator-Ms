package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Toll-Calculator/types"
)

func main() {
	listenAddr := ":30000"
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	http.HandleFunc("/aggregate", handleAggregate(svc))
	fmt.Println("HTTP transport running on port", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		err := json.NewDecoder(r.Body).Decode(&distance)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		err = svc.AggregateDistance(distance)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
