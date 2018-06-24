package lib

import (
	."github.com/alezh/novel/system/spider"
	."github.com/alezh/novel/system/http/request"
	"time"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

func init()  {
	bqg5200.Register()
}

var urlAction = map[string]string{
	"起始页":"https://www.biquge5200.cc/modules/article/toplist.php?sort=lastupdate&page=",
}

var bqg5200 = &Spider{
	Name:"biquge5200",
	Description:"笔趣阁5200.cc",
	EnableCookie:false,
	RuleTree:&RuleTree{
		Root: func(context *Context) {
			//context.SetTimer("起始页", time.Hour*24, nil)
			//起始页
			context.Aid(map[string]interface{}{"loop": [2]int{0, 2880}, "Rule": "TopList"}, "TopList")
		},
		Trunk:map[string]*Rule{
			"TopList":{
				//ItemFields: []string{
				//	"Title",
				//	"Url",
				//	"BookSite",
				//},
				AidFunc: func(context *Context, aid map[string]interface{}) interface{} {
					//循环获取书本
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++{
						context.AddQueue(&Request{
							Url:"https://www.biquge5200.cc/modules/article/toplist.php?sort=lastupdate&page="+ strconv.Itoa(loop[0]+1),
							Rule:aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(context *Context) {
					//获取书本
					query := context.GetDom()
					id := "tbody tr:not([align=\"center\"])"
					query.Find(id).Each(func(i int, selection *goquery.Selection) {
						names := selection.Find(".odd a")
						url ,_ := names.Attr("href")
						x := &Request{
							Url:          url,
							Rule:         "Cover",
							DownloaderID: 0,
						}
						context.AddQueue(x)
					})

				},
			},
			"Cover":{
				ItemFields: []string{
					"Title",
					"Author",
					"Catalog",
					"Status",
					"Desc",
					"CoverImg",
					"NewChapter",
					"Created",
					"Updated",
				},
				ParseFunc: func(context *Context) {
					query := context.GetDom()
					Author,_ := query.Find("meta[property$=author]").Attr("content")
					Title ,_ := query.Find("meta[property$=title]").Attr("content")
					Desc ,_ := query.Find("meta[property$=description]").Attr("content")
					Status ,_ := query.Find("meta[property$=status]").Attr("content")
					CoverImg ,_ := query.Find("meta[property$=image]").Attr("content")
					list := query.Find("#list dd a")
					count := list.Length()
					var index = 0
					if count >9 {
						index = 9
					}else{
						index = count - 1
					}
					Catalog := make([]bson.ObjectId,0)
					var NewChapter string
					if count >0 {
						list.Each(func(i int, selection *goquery.Selection) {
							if index >= i {
								url ,_ := selection.Attr("href")
								Name := selection.Text()
								Id := bson.NewObjectId()
								Catalog = append(Catalog,Id)
								context.AddQueue(&Request{
									Url:          url,
									Rule:         "Chapter",
									DownloaderID: 0,
									Temp: map[string]interface{}{
										"name":Name,
										"ids":Id,
									},
								})
							}else if i==0{
								fmt.Println(selection.Text())
								NewChapter = selection.Text()
							}
						})
						context.Output(map[int]interface{}{
							0: Title,
							1: Author,
							2: Catalog,
							3: Status,
							4: Desc,
							5: CoverImg,
							6: NewChapter,
							7:time.Now().Unix(),
							8:time.Now().Unix(),
						})
					}else {
						fmt.Println("url:",query.Url.Path)
					}
				},
			},
			"Chapter":{
				ItemFields:[]string{
					"_id",
					"Content",
					"Title",
				},
				ParseFunc: func(context *Context) {
					var title = context.GetTemp("name", "a").(string)
					var id = context.GetTemp("ids", "a").(bson.ObjectId)
					query := context.GetDom()
					Content := strings.Replace(strings.TrimSpace(query.Find("#content").Text()), "\n\n    ", "\n", -1)
					context.Output(map[int]interface{}{
						0:id,
						1:Content,
						2:title,
					})
				},
			},
		},
	},

}
