package shortener

import (
	"github.com/vstdy0/go-project/storage/inmemory"
	"strconv"
	"sync"
)

var mutex = sync.RWMutex{}

var urlModel = inmemory.URLModel{}
var id int

func AddURL(url string) string {
	mutex.Lock()
	id++
	urlID := strconv.Itoa(id)
	urlModel.Set(urlID, url)
	mutex.Unlock()
	return urlID
}

func GetURL(id string) string {
	return urlModel.Get(id)
}
