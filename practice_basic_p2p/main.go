package main

import (
	"net"
	"github.com/urfave/cli"
	"os"
	"log"
	"fmt"
				"time"
	"sync"
	)

var Blocks = Block {}
var mutex = &sync.Mutex{}
var isChange bool
var nodeList NodeList

func StartServer(tcpPort string, restPort string){
	nodeIP := fmt.Sprintf("%s",GetOutboundIP())
	log.Printf("[Start] IP:PORT : %s:%s", nodeIP, tcpPort)
	semiPort := ":" + tcpPort
	ln, err := net.Listen("tcp", semiPort)
	if err != nil {
		log.Panic(err)
	}

	nodeList.Identify.IP = nodeIP
	nodeList.Identify.Port = tcpPort

	go func() {
		restAPI := RestAPI{}
		restAPI.handleRequest(restPort)
	}()

	defer ln.Close()


	go func() {
		for {
			time.Sleep(3 * time.Second)
			if isChange {
				mutex.Lock()
				sendNode()
				isChange = false
				mutex.Unlock()
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleNodeExchange(conn)
		//conn.Close()
	}
}


func main() {
	app := cli.NewApp()
	app.Name = "CONCURRENCY"
	app.Usage = "PRACTICE"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "tcpPort, tp",
			Value: "",
			Usage: "set node tcpPort",
		},
		cli.StringFlag{
			Name:  "restPort, rp",
			Value: "",
			Usage: "set node restPort",
		},
	}
	app.Commands = []cli.Command{}
	app.Before = func(c *cli.Context) error {
		tcpPort := c.String("tcpPort")
		restPort := c.String("restPort")

		if tcpPort != "" && restPort != "" {
			StartServer(tcpPort, restPort)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}