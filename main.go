package main

import (
	"fmt"
	_ "github.com/alezh/gether/lib"
	"github.com/alezh/gether/system/pool"
)

func main() {
	req :=pool.GetRequest("https://www.biquge5200.cc","GET","")
	doc,_ := req.DownLoad()	
	fmt.Println(doc.Text())
}
