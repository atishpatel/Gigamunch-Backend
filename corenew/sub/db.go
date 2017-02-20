package sub

import (
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
