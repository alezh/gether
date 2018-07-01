package main

import (
	"sync"
	"runtime"
	"fmt"
	// _ "github.com/alezh/gether/lib"
	"github.com/alezh/gether/system/pool"

)

var ThreadWait  sync.WaitGroup

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ThreadWait.Add(1)
	go receiving()
	pool.Queue.InsertQueue("https://www.biquge5200.cc","GET","sx")
	
	ThreadWait.Wait()
}

func receiving(){
	var f func(map[string]*pool.Response)
	f = func(m map[string]*pool.Response){
		for v,k := range m{
			switch v {
			case "sx":
				fmt.Println(k.Node)
				break
			}
		}
	}
	for {
		value := <- pool.SuccessChan
		f(value)
   }
}