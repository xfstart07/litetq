package kvclient

import "github.com/xfstart07/litetq"

type KVClient interface {
	All(table string) ([]*litetq.Message, error)
	Get(key string) (*litetq.Message, error)
	Set(key string, msg *litetq.Message) error
	Del(key string) error
}
