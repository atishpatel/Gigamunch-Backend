package common

import "context"

// Key is a key for entry in DB.
type Key interface {
	DSKey() interface{}
	IntID() int64
	NameID() string
}

// DB is an interface for the NoSQL db. This allows clients for
// App Engine Standard, App Engine Flex, Container Engine, and Compute Engine to share code.
type DB interface {
	ErrNoSuchEntity() error
	IDKey(ctx context.Context, kind string, intID int64) Key
	NameKey(ctx context.Context, kind, name string) Key
	IncompleteKey(ctx context.Context, kind string) Key
	Put(ctx context.Context, key Key, src interface{}) (Key, error)
	PutMulti(ctx context.Context, keys []Key, src interface{}) ([]Key, error)
	Get(ctx context.Context, key Key, dst interface{}) error
	GetMulti(ctx context.Context, keys []Key, dst interface{}) error
	Query(ctx context.Context, kind string, offset, limit int, orderFieldName string, dst interface{}) ([]Key, error)
	QueryFilter(ctx context.Context, kind string, offset, limit int, filterString string, filterValue interface{}, dst interface{}) ([]Key, error)
	QueryFilterOrdered(ctx context.Context, kind string, offset, limit int, orderFieldName string, filterString string, filterValue interface{}, dst interface{}) ([]Key, error)
}
