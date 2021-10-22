package cache

import (
	"WebApplication/internal/logger"
	"log"
	"strings"
	"time"
)

const (
	DB1 = iota
	DB2
)

type Cache interface {
	get(string, int) (interface{}, error)
	set(string, interface{}, int) error
	delete(string, int) error
	clearAll()
	close(int)
	bind(interface{}) interface{}
	unbind(interface{}, int) interface{}
}

var cache Cache

func Init() {
	err := initRedis()
	if err != nil {
		logger.Log.Println(err)
		initRistRetto()
		handleError(err)
	} else {
		cache = redisDb
	}
}
func Get(key string, kind int) (value interface{}, found bool) {
	value, err := cache.get(key, kind)
	defer func() { value = cache.unbind(value, kind) }()
	if err != nil && value == "" {
		handleError(err)
		return nil, false
	}
	return value, true
}
func Set(key string, data interface{}, kind int) {
	err := cache.set(key, cache.bind(data), kind)
	if err != nil {
		logger.Log.Println(err)
		handleError(err)
	}
}
func Delete(key string, kind int) {
	err := cache.delete(key, kind)
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	if isAvailable == true && strings.Contains(err.Error(), "connect: connection refused") {
		isAvailable = false
		cache = ristdb
		t := time.NewTicker(5 * time.Second)
		go checkAvailability(t)
	}
}

func checkAvailability(ticker *time.Ticker) {
	for {
		<-ticker.C
		str, err := redisDb.rdb[DB1].Ping().Result()
		if err == nil && strings.Compare(str, "PONG") == 0 {
			log.Println("redis connection established and ready to use...")
			isAvailable = true
			cache.clearAll()
			cache = redisDb
			cache.clearAll()
			break
		}
	}
}
func Clear() {
	cache.clearAll()
}
func Close() {

}
