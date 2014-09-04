package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"time"

	"github.com/glycerine/go-tigertonic"

	"net"
	"net/http"
)

func (s *WebServer) Stop() {
	close(s.RequestStop)
	s.Tts.Close()
	<-s.Done
	WaitUntilServerDown(s.Addr)
}

func (s *WebServer) IsStopRequested() bool {
	select {
	case <-s.RequestStop:
		return true
	default:
		return false
	}
}

func WaitUntilServerUp(addr string) {
	attempt := 1
	for {
		if PortIsBound(addr) {
			return
		}
		time.Sleep(50 * time.Millisecond)
		attempt++
		if attempt > 40 {
			panic(fmt.Sprintf("could not connect to server at '%s' after 40 tries of 50msec", addr))
		}
	}
}

func WaitUntilServerDown(addr string) {
	attempt := 1
	for {
		if !PortIsBound(addr) {
			return
		}
		//fmt.Printf("WaitUntilServerUp: on attempt %d, sleep then try again\n", attempt)
		time.Sleep(50 * time.Millisecond)
		attempt++
		if attempt > 40 {
			panic(fmt.Sprintf("could always connect to server at '%s' after 40 tries of 50msec", addr))
		}
	}
}

func PortIsBound(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func NewWebServer(addr string) *WebServer {

	s := &WebServer{
		Addr:        addr,
		ServerReady: make(chan bool),
		RequestStop: make(chan bool),
		Done:        make(chan bool),
		StopSigCh:   make(chan os.Signal),
	}
	//	s.Tts = tigertonic.NewServer(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		s.ServeHTTP(w, r)
	//	}))

	s.Tts = tigertonic.NewServer(addr, http.DefaultServeMux) // nil => supply debug/pprof diagnostics

	return s
}

type WebServer struct {
	Addr        string
	ServerReady chan bool      // closed once server is listening on Addr
	RequestStop chan bool      // close this to tell server to shutdown
	Done        chan bool      // recv on this to know that server is indeed shutdown
	StopSigCh   chan os.Signal // signals will send on this to request stop
	LastReqBody string
	Tts         *tigertonic.Server
}

func (webserv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("resultsHandler (on %p) for %s running ...\n", webserv, webserv.Addr)
	fmt.Printf("resultsHandler: address '%s' is bound: %v\n", webserv.Addr, PortIsBound(webserv.Addr))
	fmt.Printf("resultsHandler: stop already request: %v\n", webserv.IsStopRequested())

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	bodyAsString := string(buf.Bytes())
	fmt.Fprintf(w, "server got request body: '%s'\n", bodyAsString)
	fmt.Printf("server %p has bodyAsString: = %s\n", webserv, bodyAsString)

	webserv.LastReqBody = bodyAsString
}

func (s *WebServer) Start() *WebServer {

	go func() {
		err := s.Tts.ListenAndServe()
		if nil != err {
			//log.Println(err) // accept tcp 127.0.0.1:3000: use of closed network connection
		}
		close(s.Done)
	}()

	/*
		signal.Notify(s.StopSigCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		go func() {
			log.Println(<-s.StopSigCh)
			s.Tts.Close()
		}()
	*/

	WaitUntilServerUp(s.Addr)
	close(s.ServerReady)
	return s
}
