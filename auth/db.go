package auth

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func addUserToken(ctx context.Context, authToken *types.AuthToken, userSessions *UserSessions) error {
	key := datastore.NewKey(ctx, KindUserSessions, authToken.User.UserID, 0, nil)
	userSessions.TokenIDs = append(userSessions.TokenIDs, TokenID{JTI: authToken.JTI, Expire: authToken.Expire})
	_, err := datastore.Put(ctx, key, userSessions)
	return err
}

func removeUserToken(ctx context.Context, authToken *types.AuthToken) error {
	userKey := datastore.NewKey(ctx, KindUserSessions, authToken.User.UserID, 0, nil)
	userSessions := &UserSessions{}
	err := datastore.Get(ctx, userKey, userSessions)
	if err != nil {
		return err
	}
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if int32(userSessions.TokenIDs[i].JTI) == authToken.JTI {
			// UserSession token should be removed
			userSessions.TokenIDs = append(userSessions.TokenIDs[:i], userSessions.TokenIDs[i+1:]...)
			break
		}
	}
	_, err = datastore.Put(ctx, userKey, userSessions)
	return err
}

// TODO test!!!
func updateToken(ctx context.Context, authToken *types.AuthToken) error {
	userKey := datastore.NewKey(ctx, KindUserSessions, authToken.User.UserID, 0, nil)
	userSessions := &UserSessions{}
	err := datastore.Get(ctx, userKey, userSessions)
	if err != nil {
		return err
	}
	needPut := false
	found := false
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if int32(userSessions.TokenIDs[i].JTI) == authToken.JTI {
			found = true
			if time.Now().Sub(userSessions.TokenIDs[i].Expire) < 30*24*time.Hour {
				// should update jti
				if userSessions.TokenIDs[i].UpdatedToJTI != 0 {
					authToken.JTI = userSessions.TokenIDs[i].UpdatedToJTI
					authToken.Expire = userSessions.TokenIDs[i].UpdateToExpire
				} else {
					newJTI := getNewJTI()
					authToken.JTI = newJTI
					userSessions.TokenIDs[i].UpdatedToJTI = newJTI
					expTime := getExpTime()
					authToken.Expire = expTime
					userSessions.TokenIDs[i].UpdateToExpire = expTime
					userSessions.TokenIDs = append(userSessions.TokenIDs, TokenID{JTI: newJTI, Expire: expTime})
					needPut = true
				}
			}
			authToken.ITA = getITATime()
			break
		}
	}
	for i := len(userSessions.TokenIDs) - 1; i >= 0; i-- {
		if userSessions.TokenIDs[i].Expire.After(time.Now()) {
			// UserSession token should be removed
			userSessions.TokenIDs = append(userSessions.TokenIDs[:i], userSessions.TokenIDs[i+1:]...)
			needPut = true
		}
	}
	if needPut {
		_, err = datastore.Put(ctx, userKey, userSessions)
		if err != nil {
			return err
		}
	}
	if !found {
		return errors.ErrInvalidToken.WithArgs("JTI session not found")
	}
	authToken.User = userSessions.User
	return nil
}

func getUserSessions(ctx context.Context, userID string, userSessions *UserSessions) error {
	userKey := datastore.NewKey(ctx, KindUserSessions, userID, 0, nil)
	return datastore.Get(ctx, userKey, userSessions)
}
