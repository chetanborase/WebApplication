package cache

import (
	"errors"
	"github.com/dgraph-io/ristretto"
)

type rist struct {
	db []*ristretto.Cache
}

var ristdb rist

func initRistRetto() {
	ristdb = rist{}
	var err error
	single, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	bulk, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ristdb.db = append(ristdb.db, single)
	ristdb.db = append(ristdb.db, bulk)
}

func (db rist) get(key string, i int) (interface{}, error) {
	val, found := db.db[i].Get(key)
	if found {
		return val, nil
	}
	return nil, errors.New("not found")
}
func (db rist) set(key string, value interface{}, i int) error {
	_ = db.db[i].Set(key, value, 1)
	if i == DB1 {
		db.clear(i)
	}
	return nil
}
func (db rist) delete(key string, i int) error {
	db.db[i].Del(key)
	return nil
}
func (db rist) close(i int) {
	db.db[i].Close()
}

func (db rist) clear(i int) {
	db.db[i].Clear()
}

func (db rist) clearAll() {
	for _, i := range db.db {
		i.Clear()
	}
}
func (db rist) bind(val interface{}) interface{} {
	return val
}
func (db rist) unbind(val interface{}, kind int) interface{} {
	return val
}
