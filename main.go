package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	// minwait is the minimum seconds to wait before interrupting
	minwait = flag.Int("minwait", 1, "minimum seconds to wait before interrupting")

	// maxwait is the maximum seconds to wait before interrupting
	maxwait = flag.Int("maxwait", 60, "maximum seconds to wait before interrupting")

	// times is the number of times to re-launch the cmd
	times = flag.Int("times", 10, "number of times to re-launch the cmd")
)

func init() {
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup

	// filter out flags from args
	// any leading arg starting with "-" is considered a flag
	i := 1
	flags := os.Args[1:]
	for _, f := range flags {
		if strings.HasPrefix(f, "-") {
			i++
		} else {
			break
		}
	}
	cmd := os.Args[i]
	args := os.Args[i+1:]

	for i := 0; i < *times; i++ {
		wg.Add(1)
		go Launch(cmd, args, i, &wg)
	}
	wg.Wait()
}

// Launch starts the passed exec cmd, sets a random timeout
// to interrupt and waits for the process to finish
func Launch(cmd string, args []string, i int, wg *sync.WaitGroup) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		log.Fatalf("%v", err)
	}

	// generate random time within minwait to maxwait
	wait := rand.Int()%(*maxwait-*minwait+1) + *minwait
	time.AfterFunc(time.Duration(wait)*time.Second, func() {
		// after random waiting time, send a SIGINT
		if err := c.Process.Signal(os.Interrupt); err != nil {
			log.Printf("%v", err)
		}
	})

	// wait for the process to finish
	// if the process has finished before timeout, allow this
	// goroutine to exit
	if err := c.Wait(); err != nil {
		log.Printf("process %v with timeout %v: %v", i, wait, err)
	}
	wg.Done()
}
