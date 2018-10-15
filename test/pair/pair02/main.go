//rhu@HZHL4 ~/code/go
//$ go install github.com/rhu1/scribble-go-runtime/test/pair/pair02
//$ bin/pair02.exe

//go:generate scribblec-param.sh Pair2.scr -d . -param Proto1 github.com/rhu1/scribble-go-runtime/test/pair/pair02/Pair2 -param-api S -param-api W

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	//"strconv"
	"sync"
	"time"

	"github.com/rhu1/scribble-go-runtime/runtime/twodim/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/rhu1/scribble-go-runtime/test/pair/pair02/Pair2/Proto1"
	S11 "github.com/rhu1/scribble-go-runtime/test/pair/pair02/Pair2/Proto1/S_l1r1tol1r1"
	W11_K "github.com/rhu1/scribble-go-runtime/test/pair/pair02/Pair2/Proto1/W_l1r1toK"
	"github.com/rhu1/scribble-go-runtime/test/util"
)

var _ = gob.Register
var _ = shm.Dial
var _ = tcp.Dial


/*
var LISTEN = tcp.Listen
var DIAL = tcp.Dial
var FORMATTER = func() *session2.GobFormatter { return new(session2.GobFormatter) } 
/*/
var LISTEN = shm.Listen
var DIAL = shm.Dial
var FORMATTER = func() *session2.PassByPointer { return new(session2.PassByPointer) } 
//*/


const PORT = 8888


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//K := session2.XY(1, 3)
	//K := session2.XY(3, 1)
	K := session2.XY(2, 3)

	wg := new(sync.WaitGroup)
	wg.Add(K.Flatten(K) + 1)

	go server_S11(wg, K)

	time.Sleep(100 * time.Millisecond) //2017/12/11 11:21:40 cannot connect to 127.0.0.1:8891: dial tcp 127.0.0.1:8891: connectex: No connection could be made because the target machine actively refused it.

	for j := session2.XY(1, 1); j.Lte(K); j = j.Inc(K) {
		go clientCode(wg, K, j)
	}

	wg.Wait()
}


func server_S11(wg *sync.WaitGroup, K session2.Pair) *S11.End {
	var err error
	P1 := Proto1.New()
	self := session2.XY(1, 1)
	S := P1.New_S_l1r1tol1r1(K, self)
	as := make([]transport2.ScribListener, K.Flatten(K))
	for j := (session2.XY(1, 1)); j.Lte(K); j = j.Inc(K) {
		as[j.Flatten(K)-1], err = LISTEN(PORT+j.Flatten(K))
		if err != nil {
			panic(err)
		}
		defer as[j.Flatten(K)-1].Close()
	}
	for j := (session2.XY(1, 1)); j.Lte(K); j = j.Inc(K) {
		err = S.W_l1r1toK_Accept(j, as[j.Flatten(K)-1], FORMATTER())
		if err != nil {
			panic(err)
		}
	}
	//fmt.Println("S ready to run")
	end := S.Run(runS)
	wg.Done()
	return &end
}


func runS(s *S11.Init) S11.End {
	data := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	end := s.W_l1r1toK_Scatter_Foo(data)
	fmt.Println("S scattered:", data)
	return *end
}

func clientCode(wg *sync.WaitGroup, K session2.Pair, self session2.Pair) *W11_K.End {
	P1 := Proto1.New()
	W := P1.New_W_l1r1toK(K, self)
	err := W.S_l1r1tol1r1_Dial(session2.XY(1,1), util.LOCALHOST, PORT+self.Flatten(K), DIAL, FORMATTER())
	if err != nil {
		panic(err)
	}
	end := W.Run(runW)
	wg.Done()
	return &end
}

func runW(w *W11_K.Init) W11_K.End {
	data := make([]int, 1)
	end := w.S_l1r1_Gather_Foo(data) // FIXME: panic: interface conversion: interface {} is int, not *int -- cf. gob.Register in commented init() ?
	fmt.Println("W("+w.Ept.Self.String()+") gathered:", data)
	return *end
}
