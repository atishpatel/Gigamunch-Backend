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
	subs, err := subC.GetHasSubscribed(time.Now())
	if err != nil {
		log.Errorf(ctx, "failed to SendStatsSMS: failed to sub.GetHasSubscribed: %s", err)
		return nil
	}
	date := time.Now()
	for date.Weekday() != time.Saturday {
		date = date.Add(-1 * time.Hour)
	}

	totalSubs := -1
	// week
	dateMinus7Days := date.Add(-1 * time.Hour * 24 * 7)
	newSubsLastWeek := 0
	cancelsLastWeek := 0
	cancelsLessThan4WeekRetention := 0
	cancels4To8WeekRetention := 0
	cancelsMoreThan8WeekRetention := 0
	totalSubsLastWeek := 0
	sumDaysWithUs := 0
	// month
	dateMinus30Days := date.Add(-1 * time.Hour * 24 * 30)
	newSubs30Days := 0
	cancelsLast30Days := 0
	totalSubs30DaysAgo := 0
	for _, sub := range subs {
		if sub.IsSubscribed {
			totalSubs++
		}
		// week
		if sub.SubscriptionDate.After(dateMinus7Days) && sub.SubscriptionDate.Before(date) {
			newSubsLastWeek++
		}
		if sub.UnSubscribedDate.After(dateMinus7Days) && sub.UnSubscribedDate.Before(date) {
			cancelsLastWeek++
			daysWithUs := int(sub.UnSubscribedDate.Sub(sub.SubscriptionDate) / (time.Hour * 24))
			sumDaysWithUs += daysWithUs
			if daysWithUs > 7*8 {
				cancelsMoreThan8WeekRetention++
			} else if daysWithUs < 7*4 {
				cancelsLessThan4WeekRetention++
			} else {
				cancels4To8WeekRetention++
			}
		}
		if sub.SubscriptionDate.Before(dateMinus7Days) {
			totalSubsLastWeek++
		}
		// month
		if sub.SubscriptionDate.After(dateMinus30Days) && sub.SubscriptionDate.Before(date) {
			newSubs30Days++
		}
		if sub.UnSubscribedDate.After(dateMinus30Days) && sub.UnSubscribedDate.Before(date) {
			cancelsLast30Days++
		}
		if sub.SubscriptionDate.Before(dateMinus30Days) {
			totalSubs30DaysAgo++
		}
	}
	weeklyChurn := (float32(cancelsLastWeek) / float32(totalSubsLastWeek)) * 100
	monthlyChurn := (float32(cancelsLast30Days) / float32(totalSubs30DaysAgo)) * 100
	avgWeeksWithUs := (float32(sumDaysWithUs) / float32(cancelsLastWeek)) / 7
	msg := `📈 %s stats: 
	Total Subs: %d

🐉
Stats for last week:
	New Subs:       %d
	Cancel Subs:   %d
	Weekly Churn: %.2f %%
	
🚧
Number of weeks with us:
	>8 weeks:    %d
	4-8 weeks:  %d
	<4 weeks:    %d
	Avg weeks: %.2f
	
🗓️
Stats for last 30 days:
	New Subs:          %d
	Cancel Subs:      %d
	Monthly Churn:  %.2f %%
	`
	msg = fmt.Sprintf(msg, date.Add(time.Hour*25).Format("Jan 2"), totalSubs, newSubsLastWeek, cancelsLastWeek, weeklyChurn, cancelsMoreThan8WeekRetention, cancels4To8WeekRetention, cancelsLessThan4WeekRetention, avgWeeksWithUs, newSubs30Days, cancelsLast30Days, monthlyChurn)
	messageC := message.New(ctx)
	numbers := []string{"6155454989", "9316445311", "9316446755", "6153975516"}
	// numbers := []string{"9316445311"}
	for _, number := range numbers {
		err = messageC.SendDeliverySMS(number, msg)
		if err != nil {
			log.Errorf(ctx, "failed to send quantity sms: %+v", err)
		}
	}
	return nil
}
