package inquiry

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64) (*Inquiry, error) {
	inquiry := new(Inquiry)
	key := datastore.NewKey(ctx, kindInquiry, "", id, nil)
	err := datastore.Get(ctx, key, inquiry)
	inquiry.ID = id
	return inquiry, err
}

func getMulti(ctx context.Context, ids []int64) ([]Inquiry, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids cannot be 0 for getMulti")
	}
	dst := make([]Inquiry, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindInquiry, "", ids[i], nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	for i := range dst {
		dst[i].ID = ids[i]
	}
	return dst, nil
}

func put(ctx context.Context, id int64, inquiry *Inquiry) error {
	var err error
	inquiry.ID = id
	key := datastore.NewKey(ctx, kindInquiry, "", id, nil)
	_, err = datastore.Put(ctx, key, inquiry)
	return err
}

func putIncomplete(ctx context.Context, inquiry *Inquiry) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindInquiry, nil)
	key, err = datastore.Put(ctx, key, inquiry)
	if key != nil {
		inquiry.ID = key.IntID()
	}
	return inquiry.ID, err
}

// getCookInquirys returns a list of Inquirys for a Cook.
func getCookInquirys(ctx context.Context, cookID string, startLimit, endLimit int) ([]*Inquiry, error) {
	return getXIDInquiriesOrderByCreatedDateTime(ctx, "CookID", cookID, startLimit, endLimit)
}

// getEaterInquirys returns a list of Inquirys for a Eater.
func getEaterInquirys(ctx context.Context, eaterID string, startLimit, endLimit int) ([]*Inquiry, error) {
	return getXIDInquiriesOrderByCreatedDateTime(ctx, "EaterID", eaterID, startLimit, endLimit)
}

func getXIDInquiriesOrderByCreatedDateTime(ctx context.Context, filterIDString string, id string, startLimit, endLimit int) ([]*Inquiry, error) {
	offset := startLimit
	limit := endLimit - startLimit
	query := datastore.NewQuery(kindInquiry).
		Filter(filterIDString+" =", id).
		Order("-CreatedDateTime").
		Offset(offset).
		Limit(limit)
	var dst []*Inquiry
	keys, err := query.GetAll(ctx, &dst)
	if err != nil {
		return nil, err
	}
	for i := range dst {
		dst[i].ID = keys[i].IntID()
	}
	return dst, nil
}

// getByBTTransactionID returns a list of Inquirys by BTTransactionID.
func getByBTTransactionID(ctx context.Context, btTransactionID string) (*Inquiry, error) {
	query := datastore.NewQuery(kindInquiry).Filter("BTTransactionID =", btTransactionID)
	var dst []Inquiry
	keys, err := query.GetAll(ctx, &dst)
	if err != nil {
		return nil, err
	}
	if len(keys) != 1 || len(dst) != 1 {
		return nil, fmt.Errorf("BTTransactionID(%s) doesn't have only one order", btTransactionID)
	}
	dst[0].ID = keys[0].IntID()
	return &dst[0], nil
}
