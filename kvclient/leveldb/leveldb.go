package leveldb

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/xfstart07/litetq"
	"github.com/xfstart07/litetq/kvclient"
)

var _ kvclient.KVClient = (*db)(nil)

type db struct {
	client *leveldb.DB
}

func New(dir string) (*db, error) {
	var (
		client *leveldb.DB
		err    error
	)

	if client, err = leveldb.OpenFile(dir, nil); err != nil {
		if client, err = leveldb.RecoverFile(dir, nil); err != nil {
			return nil, err
		}
	}

	return &db{client: client}, nil
}

// All implements kvclient.KVClient
func (db *db) All(table string) ([]*litetq.Message, error) {
	iter := db.client.NewIterator(util.BytesPrefix([]byte(table)), nil)
	defer iter.Release()

	list := make([]*litetq.Message, 0)
	for iter.Next() {
		var m litetq.Message
		if err := json.Unmarshal(iter.Value(), &m); err != nil {
			return nil, err
		}
		list = append(list, &m)
	}
	return list, nil
}

// Del implements kvclient.KVClient
func (db *db) Del(key string) error {
	return db.client.Delete([]byte(key), nil)
}

// Get implements kvclient.KVClient
func (db *db) Get(key string) (*litetq.Message, error) {
	b, err := db.client.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var m litetq.Message
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Set implements kvclient.KVClient
func (db *db) Set(key string, msg *litetq.Message) error {
	b, _ := json.Marshal(msg)

	tx := new(leveldb.Batch)
	defer tx.Reset()

	tx.Put([]byte(key), b)

	return db.client.Write(tx, nil)
}
