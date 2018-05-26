package spider

import (
	"fmt"
	"github.com/alezh/gether/system/pinyin"
)

//存放蜘蛛规则集合
type Group struct {
	list   []*Spider
	hash   map[string]*Spider
	sorted bool
}
var SpiderGroup = &Group{
	list:make([]*Spider,0),
	hash:make(map[string]*Spider,0),
}

//载入规则
func (s *Group)Load(spider *Spider) *Spider {
	name := spider.Name
	for i := 2; true; i++ {
		if _, ok := s.hash[name]; !ok {
			s.hash[name] = spider
			break
		}
		name = fmt.Sprintf("%s(%d)", spider.Name, i)
	}
	spider.Name = name
	s.list = append(s.list,spider)
	return spider
}

//取出所有规则
func (self *Group) Get () []*Spider {
	//规则按照首字母排序
	if !self.sorted {
		l := len(self.list)
		initials := make([]string, l)
		newlist := map[string]*Spider{}
		for i := 0; i < l; i++ {
			initials[i] = self.list[i].GetName()
			newlist[initials[i]] = self.list[i]
		}
		pinyin.SortInitials(initials)
		for i := 0; i < l; i++ {
			self.list[i] = newlist[initials[i]]
		}
		self.sorted = true
	}
	return self.list
}

//通过名字获取规则
func (self *Group) GetByName(name string) *Spider {
	return self.hash[name]
}