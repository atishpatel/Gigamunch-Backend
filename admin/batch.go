package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// UpdatePhoneNumbers updates phonenumbers for subscribers.
func (s *server) UpdatePhoneNumbers(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {

	return nil
}

// MigrateToNewSubscribersStruct migrates subscribers to new struct.
func (s *server) MigrateToNewSubscribersStruct(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var start, limit int64
	startStr := r.URL.Query().Get("start")
	if startStr != "" {
		start, _ = strconv.ParseInt(startStr, 10, 64)
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	log.Infof(ctx, "start: %d limit: %d", start, limit)
	suboldC := subold.NewWithLogging(ctx, log)
	err := suboldC.BatchSubscriptionSignUpToSubscriber(start, limit)
	if err != nil {
		return errors.Annotate(err, "failed to sub.BatchSubscriptionSignUpToSubscriber")
	}
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	subs, err := subC.GetActive(0, 1000)
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetActive")
	}
	_, err = subold.Subscriberget(ctx, "atish@gigamunchapp.com")
	if err != nil {
		log.Errorf(ctx, "failed to sub.Subscriberget atish@gigamunchapp.com: %+v", err)
	}
	key := s.db.IncompleteKey(ctx, "Subscriber")
	return errors.NoError.WithMessage(fmt.Sprintf("%d active Subscribers struct; key: %s", len(subs), key.NameID()))
}
