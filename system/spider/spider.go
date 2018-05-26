package spider

import (
	"sync"
)

//蜘蛛规则
type Spider struct {
	Name                   string    //名称
	Description            string    //详情
	Limit                  int64     //请求限制
	Keyin                  string
	EnableCookie           bool
	RuleTree               *RuleTree

	// 以下字段系统自动赋值
	id        int               // 自动分配的SpiderQueue中的索引
	subName   string            // 由Keyin转换为的二级标识名
	//reqMatrix *mission.Matrix   // 请求矩阵
	//timer     *Timer          // 定时器
	status    int               // 执行状态
	lock      sync.RWMutex
	once      sync.Once
}

type RuleTree struct {
	Root  func(*Context)    // 根节点(执行入口)
	Trunk map[string]*Rule // 节点散列表(执行采集过程)
}

type Rule struct {
	ItemFields []string                                              // 结果字段列表(选填，写上可保证字段顺序)
	ParseFunc  func(*Context)                                        // 内容解析函数
	AidFunc    func(*Context, map[string]interface{}) interface{} // 通用辅助函数
}

// 获取蜘蛛名称
func (self *Spider) GetName() string {
	return self.Name
}