package cache

import (
	"WebApplication/internal/storage/database/model"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

type Instance struct {
	rdb []*redis.Client
}

var (
	redisDb     Instance
	isAvailable bool
	//doOnce = &sync.Once{}
)

func initRedis() error {
	isAvailable = true
	redisDb = Instance{}
	single := redis.NewClient(&redis.Options{
		Network:    "",
		Addr:       redisUrl(),
		Dialer:     nil,
		OnConnect:  nil,
		Password:   redisPass(),
		DB:         0,
		MaxRetries: 3,
		TLSConfig:  nil,
	})
	_, err := single.Ping().Result()
	redisDb.rdb = append(redisDb.rdb, single)
	bulk := redis.NewClient(&redis.Options{
		Network:    "",
		Addr:       redisUrl(),
		Dialer:     nil,
		OnConnect:  nil,
		Password:   redisPass(),
		DB:         1,
		MaxRetries: 3,
		TLSConfig:  nil,
	})
	_, err = bulk.Ping().Result()
	redisDb.rdb = append(redisDb.rdb, bulk)
	if err != nil {
		return err
	} else {
		bulk.FlushDB()
	}
	return err
}
func (db Instance) get(key string, i int) (interface{}, error) {
	return db.rdb[i].Get(key).Result()
}
func (db Instance) set(key string, value interface{}, i int) error {
	_, err := db.rdb[i].Set(key, value, 0).Result()
	return err
}
func (db Instance) delete(key string, i int) error {
	_, err := db.rdb[i].Del(key).Result()
	return err
}
func (db Instance) close(i int) {
	_ = db.rdb[i].Close()
}
func (db Instance) clearAll() {
	for _, i := range redisDb.rdb {
		i.FlushDB()
	}
}
func (db Instance) closeAll() {
	for _, v := range db.rdb {
		v.Save()
		_ = v.Close()
	}
}
func (db Instance) onValueChange() {
	db.rdb[1].FlushDB()
}
func (db Instance) bind(val interface{}) interface{} {
	data, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (db Instance) unbind(val interface{}, kind int) interface{} {
	if val == nil {
		return nil
	}
	data := val.(string)
	if kind == DB1 {
		m := model.User{}
		_ = json.Unmarshal([]byte(data), &m)
		return m
	} else {
		m := model.SearchResult{}
		_ = json.Unmarshal([]byte(data), &m)
		return m
	}
}
