//rhu@HZHL4 ~/code/go
//$ go install github.com/rhu1/scribble-go-runtime/test/foo/foo01
//$ bin/foo01.exe

package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/rhu1/scribble-go-runtime/test/foo/foo01/Foo1/Proto1"
	S_1 "github.com/rhu1/scribble-go-runtime/test/foo/foo01/Foo1/Proto1/S_1to1"
	"github.com/rhu1/scribble-go-runtime/test/foo/foo01/Foo1/Proto1/W_1toK"
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
	/*conns :=  make([]tcp.ConnCfg, n)
	for i := 0; i < n; i++ {
		conns[i] = tcp.NewConnection("...", strconv.Itoa(PORT+i))
	}*/
	var err error
	P1 := Proto1.New()
	S := P1.New_S_1to1(K, 1)
	as := make([]*tcp.TcpListener, K)
	for j := 1; j <= K; j++ {
		as[j-1], err = tcp.Listen(PORT+j)
		if err != nil {
			panic(err)
		}
		defer as[j-1].Close()
	}
	for j := 1; j <= K; j++ {
		/*err := session.Accept(S, P1.W.Name(), i, conn)
		if err != nil {
			log.Fatalf("failed to create connection to W %d: %v", i, err)
		}*/
		S.W_1toK_Accept(j, as[j-1], new(session2.GobFormatter))
	}
	end := S.Run(runS)
	wg.Done()
	return end
}

func runS(s *S_1.Init) S_1.End {
	data := []int { 2, 3, 5, 7, 11, 13 }
	K := s.Ept.K  // Good API? -- generate param values as direct fields? (instead of generic map)
	pay := data[0:K]
	end := s.W_1toK_Scatter_A(pay)
	fmt.Println("S scattered A:", pay)
	return *end
}

func clientCode(wg *sync.WaitGroup, K int, self int) *W_1toK.End {
	P1 := Proto1.New()
	W := P1.New_W_1toK(K, self)  // Endpoint needs n to check self
	/*req := tcp.NewRequestor(util.LOCALHOST, strconv.Itoa(PORT+self))
	W.S_1to1_Dial(1, req, new(session.ScribDefaultFormatter))*/
	//W.S_1to1_Dial(1, util.LOCALHOST, strconv.Itoa(PORT+self), new(session.ScribDefaultFormatter))
	W.S_1to1_Dial(1, util.LOCALHOST, PORT+self,
		tcp.Dial,
		new(session2.GobFormatter))
	/*err := session.Connect(W, P1.S.Name(), 1, conn)
	if err != nil {
		log.Fatalf("failed to create connection to Auctioneer: %v", err)
	}*/
	end := W.Run(runW)
	wg.Done()
	return end
}

func runW(w *W_1toK.Init) W_1toK.End {
	pay := make([]int, 1)
	end := w.S_1to1_Gather_A(pay)
	fmt.Println("W(" + strconv.Itoa(w.Ept.Self) + ") gathered:", pay)
	return *end
}
