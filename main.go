package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof" // side-effect: installs handlers for /debug/pprof
)

// demonstrate web-based golang program debugging
// that is enabled/disabled by SIGUSR1/SIGUSR2

func main() {
	var web *WebServer
	addr := "localhost:6060"
	pid := os.Getpid()
	fmt.Printf("webdebug demo program started at pid: %d. Use\nkill -USR1 %d # to bring up webdebugging\nkill -USR2 %d # to shutdown webdebugging\n", pid, pid, pid)

	writePid()

	u1 := make(chan os.Signal)
	signal.Notify(u1, syscall.SIGUSR1)
	u2 := make(chan os.Signal)
	signal.Notify(u2, syscall.SIGUSR2)

	for {
		select {
		case <-u1:
			fmt.Printf("SIGUSR1 received\n")
			if web == nil {
				fmt.Printf("Starting webserver on %s\n", addr)
				web = NewWebServer(addr)
				web.Start()
			}
		case <-u2:
			fmt.Printf("SIGUSR2 received\n")
			if web != nil {
				fmt.Printf("Stopping webserver on %s\n", addr)
				web.Stop()
				web = nil
			}

		}
	}

	fmt.Printf("[done]\n")
}

func writePid() {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d\n", os.Getpid())
	ioutil.WriteFile("webdebug.pid", b.Bytes(), 0644)
}
