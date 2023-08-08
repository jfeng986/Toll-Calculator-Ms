package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	fmt.Println("HTTP transport running on port", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "obuID is missing",
			})
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "obuID is not a number",
			})
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		writeJSON(w, http.StatusOK, invoice)

		// fmt.Println(obuID)
	}
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
