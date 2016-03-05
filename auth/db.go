package auth

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func getUserSessions(ctx context.Context, userID string, userSessions *UserSessions) error {
	userKey := datastore.NewKey(ctx, kindUserSessions, userID, 0, nil)
	return datastore.Get(ctx, userKey, userSessions)
}

func putUserSessions(ctx context.Context, userID string, userSessions *UserSessions) error {
	userKey := datastore.NewKey(ctx, kindUserSessions, userID, 0, nil)
	_, err := datastore.Put(ctx, userKey, userSessions)
	return err
}
