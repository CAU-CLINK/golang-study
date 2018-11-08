package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"strings"
)

type RestAPI struct {
}

type address struct {
	Address string
}

type block struct {
	Data string
}

func nodeDiscovery(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data address

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}
	node := Node {
		IP: strings.Split(data.Address, ":")[0],
		Port: strings.Split(data.Address, ":")[1],
	}

	log.Printf("[REST] Get %s:%s", node.IP, node.Port)

	nodeList.comparison(node)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data block

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	log.Printf("[REST] Get Block %s", data.Data)

	block := Block{
		Data : data.Data,
	}

	sendBlock(block)
}

func (rest RestAPI) handleRequest(restPort string) {
	r := mux.NewRouter()
	r.HandleFunc("/peers", nodeDiscovery).Methods("POST")
	r.HandleFunc("/block", getBlock).Methods("POST")

	semiTcpPort := ":" + restPort
	if err := http.ListenAndServe(semiTcpPort, r); err != nil {
		log.Fatal(err)
	}
}