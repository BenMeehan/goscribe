package main

import (
	"flag"
	"log"
	"net"
)

const (
	HOST string = "0.0.0.0"
	PORT string = "8090"
)

func main() {
	// get ip and port from command line if provided
	host := flag.String("h", HOST, "ip address")
	port := flag.String("p", PORT, "port number")
	flag.Parse()

	// create a new broker
	arilator := newArilator()

	// start listening for commands from clients
	go arilator.start()

	// listen on the given ip and port
	listener, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalln("error starting broker :", err.Error())
	}
	defer listener.Close()

	log.Println("running broker on", *host, ":", *port)

	// endless loop to accept new clients as they come
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection :", err.Error())
			continue
		}
		// handle each client as a seperate go routine
		go arilator.newClient(conn)
	}
}
