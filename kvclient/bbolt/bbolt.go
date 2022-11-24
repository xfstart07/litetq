package bbolt

import (
	"github.com/xfstart07/litetq"
	"github.com/xfstart07/litetq/kvclient"
)

var _ kvclient.KVClient = (*db)(nil)

type db struct {
}

func (db *db) All(table string) ([]*litetq.Message, error) {
	panic("not implemented") // TODO: Implement
}

func (db *db) Get(key string) (*litetq.Message, error) {
	panic("not implemented") // TODO: Implement
}

func (db *db) Set(key string, msg *litetq.Message) error {
	panic("not implemented") // TODO: Implement
}

func (db *db) Del(key string) error {
	panic("not implemented") // TODO: Implement
}
