package storage

import (
	"fmt"
	"strconv"

	"github.com/alezh/novel/config"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var Source = NewInit()

type (
	MySqlDb interface {
		Mysql() *xorm.Engine
	}
	MgoDB interface {
		MongoDb() *MongoDb
	}
	DbConfig struct {
		DbIp    string
		DbName  string //db name
		DbPort  int    //3306
		DbUser  string //root
		DbPass  string //password
		Charset string //utf8
		Prefix  string //prefix_
		ShowSQL bool   //true则会在控制台打印出生成的SQL语句；
	}
	MgoConfig struct {
		MgoConn       string
		MgoDB         string
		PoolLimit     int
		MinPoolSize   int
		MaxIdleTimeMS int
	}
	MongoDb struct {
		Session  *mgo.Session
		Database *mgo.Database
	}
	DataSource struct {
		Mysql   *xorm.Engine
		MongoDb *MongoDb
	}
)

func NewInit() *DataSource {
	return &DataSource{
		NewMysqlSource().Mysql(),
		NewMongoDBSource().MongoDb(),
	}
}

func NewMysqlSource() MySqlDb {
	return &DbConfig{
		config.MYSQL_IP,
		config.MYSQL_DB,
		config.MYSQL_PORT,
		config.MYSQL_USER,
		config.MYSQL_PASS,
		config.CHARSET,
		config.PREFIX,
		config.SHOWSQL,
	}
}

func NewMongoDBSource() MgoDB {
	return &MgoConfig{
		config.MGO_CONN,
		config.MGO_DB,
		config.MGO_POOL,
		config.MGO_MinPoolSize,
		config.MGO_MaxIdleTimeMS,
	}
}

func (d *DbConfig) Mysql() *xorm.Engine {

	if d.DbIp == "" {
		return nil
	}

	dataSourceName := d.DbUser + ":" + d.DbPass + "@tcp(" + d.DbIp + ":" + strconv.Itoa(d.DbPort) + ")/" + d.DbName + "?charset=" + d.Charset

	engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err != nil {
		return nil
		//panic("orm failed to initialized")
		//return nil,errors.New("orm failed to initialized")
	}
	if errs := engine.Ping(); errs != nil {
		return nil
		//panic(errs.Error())
		//return nil,errors.New("orm failed to initialized")
	}
	if d.Prefix == "" {
		engine.SetTableMapper(core.SnakeMapper{})
	} else {
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, d.Prefix)
		engine.SetTableMapper(tbMapper)
	}
	//日志打印SQL
	engine.ShowSQL(d.ShowSQL)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1024)
	//设置最大打开连接数
	engine.SetMaxOpenConns(2048)

	return engine
}

func (m *MgoConfig) MongoDb() *MongoDb {
	if m.MgoConn != "" {
		//connection := fmt.Sprintf("mongodb://%s/%s?minPoolSize=%d&maxIdleTimeMS=%d", m.MgoConn, m.MgoDB,m.MinPoolSize,m.MaxIdleTimeMS)
		//fmt.Printf("mongodb://%d:$d@%s:%d,%s:%d/%d?minPoolSize=%s&maxIdleTimeMS=%s", m.MgoIP, m.MgoPort)
		//connection := "mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb?minPoolSize=0&maxIdleTimeMS=3000"
		session, err := mgo.Dial(fmt.Sprintf("mongodb://%s/%s?minPoolSize=%d&maxIdleTimeMS=%d", m.MgoConn, m.MgoDB, m.MinPoolSize, m.MaxIdleTimeMS))
		//p, err := mgop.DialStrongPool(connection, 5)
		if err != nil {
			return nil
			//panic(err.Error())
		}
		session.SetPoolLimit(m.PoolLimit)
		//session := p.AcquireSession()
		database := session.DB(m.MgoDB)
		return &MongoDb{session, database}
	}
	return &MongoDb{nil, nil}
}

//------------------------------------------------------------mgo 方法 ------------------------------------------------

func (m *MongoDb) InsetAll(table string, pdata ...interface{}) bool {
	coll := m.Database.C(table)
	err := coll.Insert(pdata...)
	if err != nil {
		fmt.Printf("InsertSync error: %v \r\n", err.Error())
		return false
	}
	return true
}

//翻页
func (m *MongoDb) PaginateNotSort(table string, search bson.M, skip, limit int, pSlice interface{}) {
	coll := m.Database.C(table)
	err := coll.Find(search).Skip(limit * (skip - 1)).Limit(limit).All(pSlice)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find table: %s  findall: %v", table, search)
		} else {
			fmt.Println(err.Error())
		}
	}
}
//数量
func (m *MongoDb) Count(table string) int {
	coll := m.Database.C(table)
	if i,err := coll.Count(); err == nil{
		return i
	}
	return 0
}

/*
=($eq)		bson.M{"name": "Jimmy Kuu"}
!=($ne)		bson.M{"name": bson.M{"$ne": "Jimmy Kuu"}}
>($gt)		bson.M{"age": bson.M{"$gt": 32}}
<($lt)		bson.M{"age": bson.M{"$lt": 32}}
>=($gte)	bson.M{"age": bson.M{"$gte": 33}}
<=($lte)	bson.M{"age": bson.M{"$lte": 31}}
in($in)		bson.M{"name": bson.M{"$in": []string{"Jimmy Kuu", "Tracy Yu"}}}
and			bson.M{"name": "Jimmy Kuu", "age": 33}
or			bson.M{"$or": []bson.M{bson.M{"name": "Jimmy Kuu"}, bson.M{"age": 31}}}
*/
func (m *MongoDb)FindAll(table string, search bson.M, pSlice interface{}) {
	coll := m.Database.C(table)
	err := coll.Find(search).All(pSlice)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find table: %s  findall: %v", table, search)
		} else {
			fmt.Println(err.Error())
		}
	}
}