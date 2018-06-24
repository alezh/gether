package config

//mysql
const (
	MYSQL_IP    string = "127.0.0.1"
	MYSQL_DB    string = "novel"
	MYSQL_USER  string = "root"
	MYSQL_PASS  string = "123456"
	MYSQL_PORT  int    = 3306
	CHARSET     string = "utf8mb4"
	PREFIX      string = "go_"
	SHOWSQL     bool   = false
)

//配置文件
const (
	CONFIG      string = "./config.ini"       // 配置文件路径
)

//mongodb
var (
	MGO_CONN     string = setting.String("mgo::connection")
	MGO_DB      string = setting.String("mgo::dbName")
	MGO_POOL    int    = setting.DefaultInt("mgo::SetPoolLimit",2048)
	MGO_MinPoolSize    int    = setting.DefaultInt("mgo::minPoolSize",0)
	MGO_MaxIdleTimeMS    int    = setting.DefaultInt("mgo::minPoolSize",2000)
)

// 运行状态
const (
	STOPPED = iota - 1
	STOP
	RUN
	PAUSE
)
// 运行模式
const (
	UNSET int = iota - 1
	OFFLINE
	SERVER
	CLIENT
)

// 数据头部信息
const (
	// 任务请求Header
	REQTASK = iota + 1
	// 任务响应流头Header
	TASK
	// 打印Header
	LOG
)

const SessionIDKey = "UserID"
