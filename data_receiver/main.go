package main

import (
	"fmt"
	"log"
	"net/http"

	"Toll-Calculator/types"

	"github.com/gorilla/websocket"
)

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
}

func (dr *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
	dr.conn = conn

	go dr.receive()
}

func (dr *DataReceiver) receive() {
	fmt.Println("new OBU connected, receiving data...")
	for {
		var data types.OBUData
		err := dr.conn.ReadJSON(&data)
		if err != nil {
			log.Println("error reading json:", err)
			continue
		}
		fmt.Printf("received OBU data from [%d] :: <lat %.2f, lon %.2f> \n", data.OBUID, data.Lat, data.Lon)
		dr.msg <- data
	}
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msg: make(chan types.OBUData, 100),
	}
}

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	fmt.Println("listening on port 30000...")
	http.ListenAndServe(":30000", nil)
}
