package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func setupBatchHandlers() {
	http.HandleFunc("/admin/batch/UpdatePhoneNumbers", handler(UpdatePhoneNumbers))
}

// UpdatePhoneNumbers updates phonenumbers for subscribers.
func UpdatePhoneNumbers(ctx context.Context, r *http.Request, log *logging.Client) Response {
	subC := sub.NewWithLogging(ctx, log)
	subs, err := subC.GetHasSubscribed(time.Now())
	if err != nil {
		return errors.Annotate(err, "failed to sub.GetHasSubscribed")
	}
	var updatedSubs []*sub.SubscriptionSignUp
	for i := range subs {
		oldNumber := subs[i].PhoneNumber
		oldRawNumber := subs[i].RawPhoneNumber
		subs[i].UpdatePhoneNumber(subs[i].PhoneNumber)
		if oldNumber != subs[i].PhoneNumber || (oldNumber != "" && oldRawNumber == "") {
			logging.Infof(ctx, "updating: %s", subs[i].Email)
			updatedSubs = append(updatedSubs, &subs[i])
		}
	}
	err = subC.Update(updatedSubs)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Update")
	}
	return errors.NoError.WithMessage(fmt.Sprintf("%d subs updated out of %d", len(updatedSubs), len(subs)))
}
