package account

import (
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func GetGigamuncher(ctx context.Context, userID string, gigamuncher *types.Gigamuncher, errChan chan<- error) {
	defer close(errChan)
	// TODO add cache stuff
	key := datastore.NewKey(ctx, types.KindGigamuncher, userID, 0, nil)
	errChan <- datastore.Get(ctx, key, gigamuncher)
}

func PutGigamuncher(ctx context.Context, userID string, gigamuncher *types.Gigamuncher, errChan chan<- error) {
	defer close(errChan)
	// TODO add cache stuff
	key := datastore.NewKey(ctx, types.KindGigamuncher, userID, 0, nil)
	var err error
	_, err = datastore.Put(ctx, key, gigamuncher)
	utils.Debugf(ctx, "in muncher put")
	errChan <- err
}
