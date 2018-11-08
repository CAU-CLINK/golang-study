package main

import (
	"log"
	"bytes"
	"encoding/gob"
	"net"
	"io"
	)


type Block struct {
	Data string
}

type Node struct {
	IP string
	Port string
}

type NodeList struct {
	Identify Node
	List	[]Node
}

func handleNodeExchange(conn net.Conn) {
	request := make([]byte, 4096)

	n, err := conn.Read(request)
	if err != nil {
		log.Panic(err)
	}

	request = request[:n]

	command := BytesToCommand(request[:12])

	switch command {
	case "addr":
		receiveNode(request)
	case "block":
		receiveBlock(request)
	}
}

func sendNode() {
	for _, node := range nodeList.List {
		address := node.IP + ":" + node.Port // address를 밑의 go func() 안에 넣으면 마지막 값이 모든 값에 적용된다. go routine 관련해서 좀 더 알아 볼 것
			go func() {
				conn, err := net.Dial("tcp", address)
				if err != nil {
					log.Panic(err)
				}

				defer conn.Close()

				payload := gobEncode(nodeList)
				request := append(CommandToBytes("addr"), payload...)
				_, err = io.Copy(conn, bytes.NewReader(request))
				if err != nil {
					log.Panic(err)
				}
			}()
	}
}

func sendBlock(block Block) {
	for _, node := range nodeList.List {
		address := node.IP + ":" + node.Port // address를 밑의 go func() 안에 넣으면 마지막 값이 모든 값에 적용된다. go routine 관련해서 좀 더 알아 볼 것
		go func() {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				log.Panic(err)
			}

			defer conn.Close()

			payload := gobEncode(block)
			request := append(CommandToBytes("block"), payload...)
			_, err = io.Copy(conn, bytes.NewReader(request))
			if err != nil {
				log.Panic(err)
			}
		}()
	}
}

func receiveNode(request []byte) {
	var buff bytes.Buffer
	var payload NodeList

	buff.Write(request[12:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if !HasNode(payload.Identify, nodeList) { // 받은 node의 정보를 list에 추가
		nodeList.List = append(nodeList.List, payload.Identify)
		log.Printf("[Node] Get %s:%s", payload.Identify.IP, payload.Identify.Port)
	}

	for _, node := range payload.List { // 받은 list랑 내 list랑 비교해서 없고, 나를 제외하고 추가
		if !HasNode(node, nodeList) && !(nodeList.Identify.IP == node.IP && nodeList.Identify.Port == node.Port){
			nodeList.List = append(nodeList.List, node)
			log.Printf("[Node] Get %s:%s", node.IP, node.Port)
			isChange = true
		}
	}
}

func receiveBlock(request []byte) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[12:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("[Block] Get Block %s", payload.Data)
}

func (nl *NodeList) comparison (node Node) {
	if !(node.IP == nl.Identify.IP && node.Port == nl.Identify.Port) {
		nl.List = append(nl.List, node)
		isChange = true
	}
}

func HasNode(s Node, elem NodeList) bool {
	for _, node := range elem.List{
		if s.IP == node.IP && s.Port == node.Port {
			return true
		}
	}

	return false
}