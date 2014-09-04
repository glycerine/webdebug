package main

import (
	"fmt"
	"log"
	"net"
)

// singleton.go : insure only one instance
// of this app is running, by binding a fixed-number port.

const MYPORT = ":10555"

// Acceptor tries to get the socket to actually bind.
// Necessary to show up in netstat and actually enforce singleton-ness.

func Acceptor(ln net.Listener) {
	for {
		ln.Accept() // this blocks until connection or error
	}
}

func bindSingletonInsurancePort() {
	port := MYPORT
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Sprintf("could not bind port '%s'; port may be in use already. ln = %v\n", port, ln))
	}
	go Acceptor(ln)
	log.Printf("got exclusive bind on port '%s'\n", port)

}

func init() {
	bindSingletonInsurancePort()
}
