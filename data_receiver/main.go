package main

import (
	"fmt"
	"log"
	"net/http"

	"Toll-Calculator/types"

	"github.com/gorilla/websocket"
)

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.wsHandler)
	log.Println("Ready to receive data from OBU clients...")
	http.ListenAndServe(":30001", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obu-data"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		log.Println("kafka producer error:", err)
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.Receive()
}

func (dr *DataReceiver) Receive() {
	log.Println("New OBU connected client connected !")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		log.Println("received message", data)
		if err := dr.produceData(data); err != nil {
			log.Println("kafka produce error:", err)
		}

	}
}
