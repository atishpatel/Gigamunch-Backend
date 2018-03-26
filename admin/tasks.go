package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func setupTasksHandlers() {
	http.HandleFunc("/admin/task/SetupTags", handler(SetupTags))
	http.HandleFunc("/admin/task/CheckPowerSensors", handler(CheckPowerSensors))
	http.HandleFunc("/admin/task/SendStatsSMS", handler(SendStatsSMS))
}

// SetupTags sets up tags for culture preview email and culture email 2 weeks in advance.
func SetupTags(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	nextCultureDate := time.Now().Add(time.Hour * 7 * 24)
	for nextCultureDate.Weekday() != time.Monday {
		nextCultureDate = nextCultureDate.Add(24 * time.Hour)
	}
	nextPreviewDate := nextCultureDate
	mailC, err := mail.NewClient(ctx, log)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:   "atish@gigamunchapp.com",
		AddTags: []mail.Tag{mail.GetCultureEmailTag(nextCultureDate), mail.GetPreviewEmailTag(nextPreviewDate)},
	}
	err = mailC.UpdateUser(mailReq)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.UpdateUser")
	}
	return errors.NoError
}

// CheckPowerSensors checks all the PowerSensors.
func CheckPowerSensors(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	healthC := healthcheck.New(ctx)
	err = healthC.CheckPowerSensors()
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to health.CheckPowerSensors")
	}
	return nil
}

// SendStatsSMS sends the stats on new and cancel sms to Chris and Piyush.
func SendStatsSMS(ctx context.Context, r *http.Request, log *logging.Client) Response {
	if !common.IsProd(projID) {
		return nil
	}
	subC := sub.New(ctx)
	subs, err := subC.GetHasSubscribed()
	if err != nil {
		log.Errorf(ctx, "failed to SendStatsSMS: failed to sub.GetHasSubscribed: %s", err)
		return nil
	}
	date := time.Now()
	for date.Weekday() != time.Saturday {
		date = date.Add(-1 * time.Hour)
	}
	dateMinus7Days := date.Add(-1 * time.Hour * 24 * 7)

	totalSubs := 0
	newSubsLastWeek := 0
	cancelsLastWeek := 0
	for _, sub := range subs {
		if sub.IsSubscribed {
			totalSubs++
			if sub.SubscriptionDate.After(dateMinus7Days) && sub.SubscriptionDate.Before(date) {
				newSubsLastWeek++
			}
		} else {
			if sub.UnSubscribedDate.After(dateMinus7Days) && sub.UnSubscribedDate.Before(date) {
				cancelsLastWeek++
			}
		}
	}
	msg := `%s stats:
	Total Subs: %d

	Stats for last week:
	New Subs: %d
	Cancel Subs: %d`
	msg = fmt.Sprintf(msg, date.Add(time.Hour*25).Format("Jan 2"), totalSubs, newSubsLastWeek, cancelsLastWeek)
	messageC := message.New(ctx)
	// "9316446755",
	numbers := []string{"6155454989"}
	for _, number := range numbers {
		err = messageC.SendDeliverySMS(number, msg)
		if err != nil {
			log.Errorf(ctx, "failed to send quantity sms: %+v", err)
		}
	}
	return nil
}
