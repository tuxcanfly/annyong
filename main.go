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
	minwait = flag.Uint("minwait", 1, "minimum seconds to wait before interrupting")

	// maxwait is the maximum seconds to wait before interrupting
	maxwait = flag.Uint("maxwait", 10, "maximum seconds to wait before interrupting")

	// times is the number of times to re-launch the cmd
	times = flag.Uint("times", 10, "number of times to re-launch the cmd")

	// parallel when true runs the cmd in parallel using goroutines
	parallel = flag.Bool("parallel", false, "when true runs the cmd in parallel using goroutines")

	// quit when true stop after receiving the first non-zero return code
	quit = flag.Bool("quit", true, "when true stop after receiving the first non-zero return code (unused if -parallel=true)")

	// verbose logging enabled
	verbose = flag.Bool("verbose", false, "whether to enable to verbose logging")
)

func init() {
	flag.Parse()
	if *minwait > *maxwait {
		log.Fatal("minwait cannot be greater than maxwait")
	}
	rand.Seed(time.Now().UnixNano())
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
	if len(os.Args) == i {
		// all args are flags
		return
	}
	cmd := os.Args[i]
	args := os.Args[i+1:]

	for i := 0; i < int(*times); i++ {
		wg.Add(1)
		if *parallel {
			go Launch(cmd, args, i, &wg)
		} else {
			ok := Launch(cmd, args, i, &wg)
			if !ok && *quit {
				return
			}
		}
	}
	wg.Wait()
}

// alog logs the given string
func alog(s string, v ...interface{}) {
	if *verbose {
		log.Printf(s, v)
	}
}

// Launch starts the passed exec cmd, sets a random timeout
// to interrupt and waits for the process to finish
func Launch(cmd string, args []string, i int, wg *sync.WaitGroup) (success bool) {
	defer wg.Done()

	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		log.Fatalf("%v", err)
	}

	// generate random time within minwait to maxwait
	wait := rand.Int()%int(*maxwait-*minwait+1) + int(*minwait)
	time.AfterFunc(time.Duration(wait)*time.Second, func() {
		// after random waiting time, send a SIGINT
		if err := c.Process.Signal(os.Interrupt); err != nil {
			alog(err.Error())
		}
	})

	// wait for the process to finish
	// if the process has finished before timeout, allow this
	// goroutine to exit
	if err := c.Wait(); err != nil {
		alog("process %v with timeout %v: %v", i, wait, err)
		return false
	}
	return true
}
