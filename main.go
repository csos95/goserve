package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	openBrowser   bool
	maxTries      int
	port          int
	directory     string
	exclude       string
	excludedFiles []string
)

func parseFlags() {
	pflag.BoolVarP(&openBrowser, "openBrowser", "o", false, "open the default browser")
	pflag.IntVarP(&maxTries, "maxTries", "m", 10, "max number of ports to try")
	pflag.IntVarP(&port, "port", "p", 8080, "port to use (tries random ports [8000,9000) if in use)")
	pflag.StringVarP(&directory, "directory", "d", ".", "directory to serve")
	pflag.StringVarP(&exclude, "exclude", "e", "", "comma separated list of files to exclude")
	pflag.Parse()

	// trim trailing '/' in path
	directory = strings.TrimRight(directory, "/")
	// split excluded files
	excludedFiles = strings.Split(exclude, ",")
	fmt.Printf("Excluding %v\n", excludedFiles)
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

//openBrowser opens the default user browser with the specified url
//taken from github.com/rodzzlessa24/openbrowser.go
func openInBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println(path)

	excluded := false
	for _, excludedFile := range excludedFiles {
		if strings.Contains(path, excludedFile) && excludedFile != "" {
			log.Printf("Matched excluded file '%s'", excludedFile)
			excluded = true
		}
	}

	if excluded {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, http.StatusForbidden)
	} else {
		http.ServeFile(w, r, path)
	}
}

func main() {
	parseFlags()

	err := findOpenPort()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Serving %s on port %d\n", directory, port)

	if openBrowser {
		go func() {
			time.Sleep(time.Second)
			openInBrowser(fmt.Sprintf("http://0.0.0.0:%d", port))
		}()
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Println(err)
	}
}
