package config

import (
	"os"
	"github.com/alezh/gether/system/config"
	"strconv"
	"fmt"
)

var (
	CRAWLS_CAP int = setting.DefaultInt("crawlcap", 50)
)

const (
	mode        int    = UNSET        // 节点角色
	port        int    = 2015         // 主节点端口
	master      string = "127.0.0.1"  // 服务器(主节点)地址，不含端口
	thread      int    = 20           // 全局最大并发量
	pause       int64  = 300          // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
	//outtype     string = "csv"        // 输出方式
	dockercap   int    = 10000        // 分段转储容器容量
	limit       int64  = 0            // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	proxyminute int64  = 0            // 代理IP更换的间隔分钟数
	success     bool   = true         // 继承历史成功记录
	failure     bool   = true         // 继承历史失败记录
)

var setting = set()

func set() config.Configer {
	iniconf, err := config.NewConfig("ini", CONFIG)
	if err != nil {
		file, err := os.Create(CONFIG)
		file.Close()
		iniconf, err = config.NewConfig("ini", CONFIG)
		if err != nil {
			panic(err)
		}
		defaultConfig(iniconf)
		iniconf.SaveConfigFile(CONFIG)
	} else {
		trySet(iniconf)
	}
	return iniconf
}

func defaultConfig(iniconf config.Configer) {

	iniconf.Set("mgo::connection", "127.0.0.1:27017")
	iniconf.Set("mgo::dbName", "BookDb")
	iniconf.Set("mgo::SetPoolLimit", "2048")
	iniconf.Set("mgo::minPoolSize", "0")
	iniconf.Set("mgo::maxIdleTimeMS", "2000")

	iniconf.Set("mysql::host", MYSQL_IP)
	iniconf.Set("mysql::dbname", MYSQL_DB)
	iniconf.Set("mysql::username", MYSQL_USER)
	iniconf.Set("mysql::password", MYSQL_PASS)
	iniconf.Set("mysql::port", strconv.Itoa(MYSQL_PORT))
	iniconf.Set("mysql::password", MYSQL_PASS)
	iniconf.Set("mysql::charset", CHARSET)
	iniconf.Set("mysql::prefix", PREFIX)
	iniconf.Set("mysql::showSql", fmt.Sprint(SHOWSQL))
	iniconf.Set("mysql::SetMaxIdleConns", "1024")
	iniconf.Set("mysql::SetMaxOpenConns", "2048")

	iniconf.Set("run::mode", strconv.Itoa(mode))
	iniconf.Set("run::port", strconv.Itoa(port))
	iniconf.Set("run::master", master)
	iniconf.Set("run::thread", strconv.Itoa(thread))
	iniconf.Set("run::pause", strconv.FormatInt(pause, 10))
	iniconf.Set("run::outtype", "mgo")
	iniconf.Set("run::dockercap", strconv.Itoa(dockercap))
	iniconf.Set("run::limit", strconv.FormatInt(limit, 10))
	iniconf.Set("run::proxyminute", strconv.FormatInt(proxyminute, 10))
	iniconf.Set("run::success", fmt.Sprint(success))
	iniconf.Set("run::failure", fmt.Sprint(failure))
}

func trySet(iniconf config.Configer) {

	if v := iniconf.String("mgo::connection"); v == "" {
		iniconf.Set("mgo::connection", "127.0.0.1:27017")
	}
	if v := iniconf.String("mgo::dbName"); v == "" {
		iniconf.Set("mgo::dbName", "BookDb")
	}
	if v , e := iniconf.Int("mgo::SetPoolLimit"); v <= 0 || e!=nil {
		iniconf.Set("mgo::SetPoolLimit", "2048")
	}
	if v , e := iniconf.Int("mgo::minPoolSize"); v < 0 || e!=nil {
		iniconf.Set("mgo::minPoolSize", "0")
	}
	if v , e := iniconf.Int("mgo::maxIdleTimeMS"); v <= 0 || e!=nil {
		iniconf.Set("mgo::maxIdleTimeMS", "2000")
	}


	if v := iniconf.String("mysql::host"); v == "" {
		iniconf.Set("mysql::host", MYSQL_IP)
	}
	if v := iniconf.String("mysql::dbname"); v == "" {
		iniconf.Set("mysql::dbname", MYSQL_DB)
	}
	if v := iniconf.String("mysql::username"); v == "" {
		iniconf.Set("mysql::username", MYSQL_USER)
	}
	if v := iniconf.String("mysql::password"); v == "" {
		iniconf.Set("mysql::password", MYSQL_PASS)
	}
	if v,e := iniconf.Int("mysql::port"); v <= 0 || e !=nil {
		iniconf.Set("mysql::port", strconv.Itoa(MYSQL_PORT))
	}
	if v := iniconf.String("mysql::charset"); v == "" {
		iniconf.Set("mysql::charset", CHARSET)
	}
	if v := iniconf.String("mysql::prefix"); v == "" {
		iniconf.Set("mysql::prefix", PREFIX)
	}
	if _, v := iniconf.Bool("mysql::showSql"); v != nil {
		iniconf.Set("mysql::showSql", fmt.Sprint(SHOWSQL))
	}
	if v,e := iniconf.Int("mysql::SetMaxIdleConns");  v <= 0 || e !=nil {
		iniconf.Set("mysql::SetMaxIdleConns", "1024")
	}
	if v,e := iniconf.Int("mysql::SetMaxOpenConns");  v <= 0 || e !=nil {
		iniconf.Set("mysql::SetMaxOpenConns", "2048")
	}


	if v, e := iniconf.Int("run::mode"); v < UNSET || v > CLIENT || e != nil {
		iniconf.Set("run::mode", strconv.Itoa(mode))
	}

	if v, e := iniconf.Int("run::port"); v <= 0 || e != nil {
		iniconf.Set("run::port", strconv.Itoa(port))
	}

	if v := iniconf.String("run::master"); v == "" {
		iniconf.Set("run::master", master)
	}

	if v, e := iniconf.Int("run::thread"); v <= 0 || e != nil {
		iniconf.Set("run::thread", strconv.Itoa(thread))
	}

	if v, e := iniconf.Int64("run::pause"); v < 0 || e != nil {
		iniconf.Set("run::pause", strconv.FormatInt(pause, 10))
	}

	if v := iniconf.String("run::outtype"); v == "" {
		iniconf.Set("run::outtype", "mgo")
	}

	if v, e := iniconf.Int("run::dockercap"); v <= 0 || e != nil {
		iniconf.Set("run::dockercap", strconv.Itoa(dockercap))
	}

	if v, e := iniconf.Int64("run::limit"); v < 0 || e != nil {
		iniconf.Set("run::limit", strconv.FormatInt(limit, 10))
	}

	if v, e := iniconf.Int64("run::proxyminute"); v <= 0 || e != nil {
		iniconf.Set("run::proxyminute", strconv.FormatInt(proxyminute, 10))
	}

	if _, e := iniconf.Bool("run::success"); e != nil {
		iniconf.Set("run::success", fmt.Sprint(success))
	}

	if _, e := iniconf.Bool("run::failure"); e != nil {
		iniconf.Set("run::failure", fmt.Sprint(failure))
	}

	iniconf.SaveConfigFile(CONFIG)
}