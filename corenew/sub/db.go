package sub

import (
	"errors"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

const (
	kindSubscriptionSignUp = "ScheduleSignUp"
)

func get(ctx context.Context, id string) (*SubscriptionSignUp, error) {
	i := new(SubscriptionSignUp)
	key := datastore.NewKey(ctx, kindSubscriptionSignUp, id, 0, nil)
	err := datastore.Get(ctx, key, i)
	return i, err
}

func getMulti(ctx context.Context, ids []string) ([]*SubscriptionSignUp, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids cannot be 0 for getMulti")
	}
	dst := make([]*SubscriptionSignUp, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		if ids[i] == "" {
			return nil, errors.New("ids cannot contain an empty string")
		}
		keys[i] = datastore.NewKey(ctx, kindSubscriptionSignUp, ids[i], 0, nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	// for i := range dst {
	// 	dst[i].ID = keys[i].StringID()
	// }
	return dst, nil
}

func put(ctx context.Context, id string, i *SubscriptionSignUp) error {
	var err error
	key := datastore.NewKey(ctx, kindSubscriptionSignUp, id, 0, nil)
	_, err = datastore.Put(ctx, key, i)
	return err
}

// getSubscribers returns the list of Subscribers for that day.
func getSubscribers(ctx context.Context, subDay string) ([]SubscriptionSignUp, error) {
	query := datastore.NewQuery(kindSubscriptionSignUp).
		Filter("IsSubscribed=", true).
		Filter("SubscriptionDay=", subDay).
		Limit(1000)
	var results []SubscriptionSignUp
	_, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}