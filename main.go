package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/exec"
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
	for i := 0; i < *times; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			cmd := exec.Command(os.Args[1], os.Args[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Fatalf("%v", err)
			}
			wait := rand.Int()%(*maxwait-*minwait+1) + *minwait
			time.AfterFunc(time.Duration(wait)*time.Second, func() {
				if err := cmd.Process.Signal(os.Interrupt); err != nil {
					log.Fatalf("%v", err)
				}
			})
			if err := cmd.Wait(); err != nil {
				log.Printf("process %v with timeout %v: %v", i, wait, err)
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
