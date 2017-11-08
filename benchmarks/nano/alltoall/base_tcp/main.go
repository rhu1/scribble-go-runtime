package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/nickng/scribble-go/runtime/transport/tcp"
)

const (
	NCPU   = 7
	NITERS = 100000
)

func Avg(d time.Duration, v int) float64 {
	return float64(d.Nanoseconds()) / float64(v)
}

var ncpu, niters int
var ncpu1, ncpu2 int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.IntVar(&ncpu, "ncpu", NCPU, "GOMAXPROCS")
	flag.IntVar(&niters, "niters", NITERS, "ITERS")
	flag.Parse()

	ncpu1 = ncpu / 2
	ncpu2 = ncpu - ncpu1
	if ncpu2 == 0 {
		ncpu2 = 1
	}
	wg := new(sync.WaitGroup)
	wg.Add(ncpu)

	conn := make([][]tcp.ConnCfg, ncpu1)
	for i := 0; i < ncpu1; i++ {
		conn[i] = make([]tcp.ConnCfg, ncpu2)
		for j := 0; j < ncpu2; j++ {
			conn[i][j] = tcp.NewConnection("127.0.0.1", strconv.Itoa(33333+i*(ncpu2)+j))
		}
	}

	serverCode := func(idx int) func() {

		cnn := make([](*tcp.Conn), ncpu2)
		rwm := new(sync.RWMutex)
		// One connection for each participant in the group
		for i := 1; i <= ncpu2; i++ {
			go func(i int) {
				rwm.Lock()
				cnn[i-1] = conn[idx][i-1].Accept().(*tcp.Conn)
				rwm.Unlock()
			}(i)
		}

		return func() {

			for i := 0; i < niters; i++ {
				rwm.RLock()
				for _, cn := range cnn {
					cn.Send(42)
				}
				rwm.RUnlock()
			}
			wg.Done()
		}
	}

	servers := make([]func(), ncpu1)
	for i := 0; i < ncpu1; i++ {
		servers[i] = serverCode(i)
	}

	for i := 1; i <= 1000000; i++ {
		// time waster instead of time.Sleep
	}

	clientCode := func(idx int) func() {
		var tmp int

		cnn := make([](*tcp.Conn), ncpu1)
		rwm := new(sync.RWMutex)
		// One connection for each participant in the group
		for i := 1; i <= ncpu1; i++ {
			rwm.Lock()
			cnn[i-1] = conn[i-1][idx].Connect().(*tcp.Conn)
			rwm.Unlock()
		}

		return func() {
			for i := 0; i < niters; i++ {
				rwm.RLock()
				for _, cn := range cnn {
					err := cn.Recv(&tmp)
					if err != nil {
						log.Fatalf("wrong value from server at %d: %s", i, err)
					}
				}
				rwm.RUnlock()
			}
			wg.Done()
		}
	}

	clients := make([]func(), ncpu2)
	for i := 0; i < ncpu2; i++ {
		clients[i] = clientCode(i)
	}

	run_startt := time.Now()
	for i := 1; i <= ncpu1; i++ {
		go servers[i-1]()
	}
	for i := 1; i <= ncpu2; i++ {
		go clients[i-1]()
	}
	wg.Wait()
	run_endt := time.Now()
	fmt.Println(Avg(run_endt.Sub(run_startt), niters))
}