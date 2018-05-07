//rhu@HZHL4 ~/code/go
//$ go install github.com/rhu1/scribble-go-runtime/test/foreach/foreach02
//$ bin/foreach02.exe

package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/rhu1/scribble-go-runtime/test/foreach/foreach02/Foreach2/Proto1"
	S_1 "github.com/rhu1/scribble-go-runtime/test/foreach/foreach02/Foreach2/Proto1/S_1to1"
	"github.com/rhu1/scribble-go-runtime/test/foreach/foreach02/Foreach2/Proto1/W_1toK"
	"github.com/rhu1/scribble-go-runtime/test/util"
)

const PORT = 8888

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	K := 3

	wg := new(sync.WaitGroup)
	wg.Add(K + 1)

	go serverCode(wg, K)

	time.Sleep(100 * time.Millisecond) //2017/12/11 11:21:40 cannot connect to 127.0.0.1:8891: dial tcp 127.0.0.1:8891: connectex: No connection could be made because the target machine actively refused it.

	for j := 1; j <= K; j++ {
		go clientCode(wg, K, j)
	}

	wg.Wait()
}

func serverCode(wg *sync.WaitGroup, K int) *S_1.End {
	var err error
	P1 := Proto1.New()
	S := P1.New_S_1to1(K, 1)
	as := make([]*tcp.TcpListener, K)
	//as := make([]*shm.ShmListener, K)
	for j := 1; j <= K; j++ {
		as[j-1], err = tcp.Listen(PORT+j)
		//as[j-1], err = shm.Listen(PORT+j)
		if err != nil {
			panic(err)
		}
		defer as[j-1].Close()
	}
	for j := 1; j <= K; j++ {
		err := S.W_1toK_Accept(j, as[j-1], 
			new(session2.GobFormatter))
			//new(session2.PassByPointer))
		if err != nil {
			panic(err)
		}
	}
	//fmt.Println("S ready to run")
	end := S.Run(runS)
	wg.Done()
	return end
}

func runS(s *S_1.Init_12) S_1.End {
	return *s.Foreach(nested)
}

func nested(s *S_1.Init_10) S_1.End {
	pay := make([]int, 1)
	end := s.W_ItoI_Gather_A(pay)
	fmt.Println("S gathered A:", pay)
	return *end
}

func clientCode(wg *sync.WaitGroup, K int, self int) *W_1toK.End {
	P1 := Proto1.New()
	W := P1.New_W_1toK(K, self)  // Endpoint needs n to check self
	err := W.S_1to1_Dial(1, util.LOCALHOST, PORT+self,
			tcp.Dial, new(session2.GobFormatter))
			//shm.Dial, new(session2.PassByPointer))
	if err != nil {
		panic(err)
	}
	//fmt.Println("W(" + strconv.Itoa(W.Self) + ") ready to run")
	end := W.Run(runW)
	wg.Done()
	return end
}

func runW(w *W_1toK.Init_4) W_1toK.End {
	data := []int { 2, 3, 5, 7, 11, 13 }
	self := w.Ept.Self  // Good API? -- generate param values as direct fields? (instead of generic map)
	pay := data[self:self+1]
	end := w.S_1to1_Scatter_A(pay)
	fmt.Println("W(" + strconv.Itoa(w.Ept.Self) + ") scattered:", pay)
	return *end
}
