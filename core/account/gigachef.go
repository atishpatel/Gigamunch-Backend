package account

import (
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func GetGigachef(ctx context.Context, userID string, gigachef *types.Gigachef, errChan chan<- error) {
	defer close(errChan)
	// TODO add cache stuff
	key := datastore.NewKey(ctx, types.KindGigachef, userID, 0, nil)
	errChan <- datastore.Get(ctx, key, gigachef)
}

func PutGigachef(ctx context.Context, userID string, gigachef *types.Gigachef, errChan chan<- error) {
	defer close(errChan)
	// TODO add cache stuff
	key := datastore.NewKey(ctx, types.KindGigachef, userID, 0, nil)
	var err error
	_, err = datastore.Put(ctx, key, gigachef)
	errChan <- err
}
