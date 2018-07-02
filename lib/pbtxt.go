package lib

import (
	"sync"

	. "github.com/alezh/gether/storage"
	"github.com/globalsign/mgo/bson"
)

type PbTxt struct {
	Web             string
	WebUrl          string
	LastUpUrl       string
	NewCreateUrl    string
	UnDesc          string
	NewBookPageSize int
	WaitGroup       *sync.WaitGroup
}

func init() {
	if count := Source.MongoDb.Count("BookCover"); count > 0 {
		bookCover := make([]BooksCache, 0)
		Source.MongoDb.FindAll("BookCover", bson.M{}, pageSize, &bookCover)
	}
}
