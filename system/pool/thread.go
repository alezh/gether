package pool

import (
	"time"
	"runtime"
	"fmt"
	"sync"
	"github.com/PuerkitoBio/goquery"
)

//页面返回数据

type Response struct {
	Node   *goquery.Document
	Method string
}

var Queue = newQueue()

var Error map[string]int

var SuccessChan chan map[string]*Response

type MQueue struct {
	TotalThread int
	mutex       sync.Mutex
}

var requestPool = &sync.Pool{
	New: func() interface{} {
		return Request{}
	},
}

//取出
func Load(url string, method string, action string,postData string) Request {
	resule := requestPool.Get().(Request)
	resule.Url = url
	resule.Method = method
	resule.PostData = postData
	resule.Action = action
	return resule
}

//存入
func PutRequest(req Request) {
	requestPool.Put(req)
}

func Get() Request {
	return requestPool.Get().(Request)
}

func newQueue() *MQueue {
	mq := &MQueue{
		TotalThread: 0,
	}
	SuccessChan = make(chan map[string]*Response)
	Error = make(map[string]int)
	go mq.timer()
	return mq
}

func (m *MQueue) InsertQueue(url string, method string,action string) *MQueue {
	m.mutex.Lock()
	m.TotalThread++
	PutRequest(Load(url, method, action, ""))
	m.mutex.Unlock()
	if m.TotalThread < 21 {		
		m.runTask(action)
	}
	return m
}

func (m *MQueue) runTask(action string) {
	m.log()
	for ;m.TotalThread > 0; {
		m.mutex.Lock()
		reqs := Get()
		m.TotalThread--
		m.mutex.Unlock()
		if doc, ok := reqs.DownLoad(); ok {
			var m = make(map[string]*Response)
			m[action] = &Response{
				Node:doc,
			}
			SuccessChan <- m
		} else {
			m.WrongInset(reqs)
		}
	}
}

func (m *MQueue) WrongInset(req Request) {
	if v, o := Error[req.Url]; o {
		if v < 6 {
			Error[req.Url] = v + 1
			PutRequest(req)
		} else {
			delete(Error, req.Url)
		}
		return
	}
}

func (x *MQueue)timer()  {
	timer := time.NewTicker(time.Second * 5)
	for{
		select {
		case <- timer.C:
			x.log()
		}
	}
}

func (m *MQueue)log(){
	green   := string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	reset   := string([]byte{27, 91, 48, 109})
	blue    := string([]byte{27, 91, 57, 55, 59, 52, 52, 109})

	// fmt.Printf("Waiting task => %s %d %s | connections => %s %d %s | NumGoroutine => %s %d %s \n",green,TotalThread,reset,blue,x.NewThread-x.OnlineTask.Count(),reset,blue,runtime.NumGoroutine(),reset)
	fmt.Printf("connections => %s %d %s | NumGoroutine => %s %d %s \n",green,m.TotalThread,reset,blue,runtime.NumGoroutine(),reset)
	// //捞出阻塞的数据
	// if m.OnlineTask.Count()>0 && x.NewThread <= 0 {
	// 	x.WaitingChan <- 1
	// }else if x.OnlineTask.Count() == 0 && x.NewThread <= 0 {
	// 	//循环3次 等待15秒
	// 	x.WaitGroup.Done()
	// }
}