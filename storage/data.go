package storage

import (
	"github.com/globalsign/mgo/bson"
)

type BooksCache struct {
	Id         bson.ObjectId   `bson:"_id" json:"_id"`
	Title      string          //书名
	Author     string          //作者
	CatalogUrl []*OriginUrl    //目录链接
	Catalog    []bson.ObjectId //目录
	Desc       string
}

type OriginUrl struct {
	Name string
	Url  string
}
