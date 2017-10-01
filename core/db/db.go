package db

// This package will make transition to flexible env easier.

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"google.golang.org/api/option"
	"google.golang.org/appengine/datastore"
)

var ()

// Key represents the datastore key for a stored entity. A Key can be based on either an IntID or NameID.
type Key struct {
	key *datastore.Key
}

// DSKey returns the Datastore Key.
func (k *Key) DSKey() interface{} {
	return k.key
}

// IntID returns the IntID of the Key.
func (k *Key) IntID() int64 {
	return k.key.IntID()
}

// NameID returns the NameID of the Key.
func (k *Key) NameID() string {
	return k.key.StringID()
}

// Client is a client for reading and writing data in a datastore dataset.
type Client struct {
	ctx       context.Context
	projectID string
}

// NewClient creates a new Client for a given dataset. If the project ID is empty, it is derived from the
// DATASTORE_PROJECT_ID environment variable. If the DATASTORE_EMULATOR_HOST environment variable is set,
// client will use its value to connect to a locally-running datastore emulator.
func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	return &Client{
		ctx:       ctx,
		projectID: projectID,
	}, nil
}

func (c *Client) ErrNoSuchEntity() error {
	return datastore.ErrNoSuchEntity
}

// IDKey creates a new key with an ID. The supplied kind cannot be empty.
// The supplied parent must either be a complete key or nil. The namespace of the new key is empty.
func (c *Client) IDKey(ctx context.Context, kind string, intID int64) common.Key {
	return &Key{key: datastore.NewKey(ctx, kind, "", intID, nil)}
}

// NameKey creates a new key with a name. The supplied kind cannot be empty.
// The supplied parent must either be a complete key or nil. The namespace of the new key is empty.
func (c *Client) NameKey(ctx context.Context, kind, name string) common.Key {
	return &Key{key: datastore.NewKey(ctx, kind, name, 0, nil)}
}

// IncompleteKey creates a new incomplete key.
// The supplied kind cannot be empty. The namespace of the new key is empty.
func (c *Client) IncompleteKey(ctx context.Context, kind string) common.Key {
	return &Key{
		key: datastore.NewIncompleteKey(ctx, kind, nil),
	}
}

// Put saves the entity src into the datastore with key k. src must be a struct
// pointer or implement PropertyLoadSaver; if a struct pointer then any
// unexported fields of that struct will be skipped. If k is an incomplete key,
// the returned key will be a unique key generated by the datastore.
func (c *Client) Put(ctx context.Context, key common.Key, src interface{}) (common.Key, error) {
	k, err := datastore.Put(ctx, key.DSKey().(*datastore.Key), src)
	if err != nil {
		return nil, err
	}
	return &Key{key: k}, nil
}

// PutMulti is a batch version of Put.
//
// src must satisfy the same conditions as the dst argument to GetMulti.
func (c *Client) PutMulti(ctx context.Context, keys []common.Key, src interface{}) ([]common.Key, error) {
	dsKeys := make([]*datastore.Key, len(keys))
	for i := range keys {
		dsKeys[i] = keys[i].DSKey().(*datastore.Key)
	}
	rDSKeys, err := datastore.PutMulti(ctx, dsKeys, src)
	if err != nil {
		return nil, err
	}
	rKeys := make([]common.Key, len(rDSKeys))
	for i := range rDSKeys {
		rKeys[i] = &Key{key: rDSKeys[i]}
	}
	return rKeys, nil
}

// Get loads the entity stored for key into dst, which must be a struct pointer or implement PropertyLoadSaver.
// If there is no such entity for the key, Get returns ErrNoSuchEntity.
//
// The values of dst's unmatched struct fields are not modified, and matching slice-typed fields are not reset
// before appending to them. In particular, it is recommended to pass a pointer to a zero valued struct on each Get call.
//
// ErrFieldMismatch is returned when a field is to be loaded into a different type than the one it was stored
// from, or when a field is missing or unexported in the destination struct. ErrFieldMismatch is only returned
// if dst is a struct pointer.
func (c *Client) Get(ctx context.Context, key common.Key, dst interface{}) error {
	err := datastore.Get(ctx, key.DSKey().(*datastore.Key), dst)
	if err == datastore.ErrNoSuchEntity {
		return c.ErrNoSuchEntity()
	}
	return err
}

// GetMulti is a batch version of Get.
//
// dst must be a []S, []*S, []I or []P, for some struct type S, some interface
// type I, or some non-interface non-pointer type P such that P or *P
// implements PropertyLoadSaver. If an []I, each element must be a valid dst
// for Get: it must be a struct pointer or implement PropertyLoadSaver.
//
// As a special case, PropertyList is an invalid type for dst, even though a
// PropertyList is a slice of structs. It is treated as invalid to avoid being
// mistakenly passed when []PropertyList was intended.
func (c *Client) GetMulti(ctx context.Context, keys []common.Key, dst interface{}) error {
	dsKeys := make([]*datastore.Key, len(keys))
	for i := range keys {
		dsKeys[i] = keys[i].DSKey().(*datastore.Key)
	}
	return datastore.GetMulti(ctx, dsKeys, dst)
}

// QueryFilter runs a query with filter parameter.
func (c *Client) QueryFilter(ctx context.Context, kind string, filterString string, filterValue interface{}, dst interface{}) ([]common.Key, error) {
	dsKeys, err := datastore.NewQuery(kind).Filter(filterString, filterValue).GetAll(ctx, dst)
	if err != nil {
		return nil, err
	}
	keys := make([]common.Key, len(dsKeys))
	for i := range dsKeys {
		keys[i] = &Key{key: dsKeys[i]}
	}
	return keys, nil
}

// Query runs a query with parameters.
func (c *Client) Query(ctx context.Context, kind string, offset, limit int, orderFieldName string, dst interface{}) ([]common.Key, error) {
	dsKeys, err := datastore.NewQuery(kind).Offset(offset).Limit(limit).Order(orderFieldName).GetAll(ctx, dst)
	if err != nil {
		return nil, err
	}
	keys := make([]common.Key, len(dsKeys))
	for i := range dsKeys {
		keys[i] = &Key{key: dsKeys[i]}
	}
	return keys, nil
}
