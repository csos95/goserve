package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	port      int
	maxTries  int
	directory string
)

func parseFlags() {
	flag.IntVar(&port, "p", 8080, "port to use (tries random ports [8000,9000) if in use)")
	flag.IntVar(&maxTries, "m", 10, "max number of ports to try")
	flag.StringVar(&directory, "d", ".", "directory to serve")
	flag.Parse()

	// trim trailing '/' in path
	directory = strings.TrimRight(directory, "/")
}

func findOpenPort() error {
	var tries int
	var err error
	rand.Seed(time.Now().Unix())

	// while err is nil, we were able to successfully
	// connect to that port; therefore, that port is in use.
	// there are probably other cases where this will
	// fail, but for now it should be good enough
	for err == nil {
		tries++
		_, err = net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err == nil {
			// the port is in use, try another random port in [8000,9000)
			port = rand.Intn(1000) + 8000
		} else if tries == maxTries {
			return fmt.Errorf("couldn't find an open port after %d tries", maxTries)
		}
	}

	return nil
}

func main() {
	parseFlags()

	err := findOpenPort()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Serving %s on port %d\n", directory, port)

	http.Handle("/", http.FileServer(http.Dir(directory)))
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Println(err)
	}
}
