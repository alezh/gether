package pool

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//页面返回数据

type Response struct {
	Node   *goquery.Document
	Method string
}

var requestPool = &sync.Pool{
	New: func() interface{} {
		return Request{}
	},
}

//取出
func GetRequest(url string, method string, postData string) Request {
	resule := requestPool.Get().(Request)
	resule.Url = url
	resule.Method = method
	resule.PostData = postData
	return resule
}

//存入
func PutRequest(req Request) {
	requestPool.Put(req)
}
