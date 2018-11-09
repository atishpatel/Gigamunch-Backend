package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// TODO: UpdateDrip

func (s *server) ProcessActivity(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	parms, err := tasks.ParseProcessSubscriptionRequest(r)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseProcessSubscriptionRequest. Err:%+v", err)
	}

	activityC, _ := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	err = activityC.Process(parms.Date, parms.SubEmail)
	if err != nil {
		log.Errorf(ctx, "failed to activity.Process(Date:%s SubEmail:%s). Err:%+v", parms.Date, parms.SubEmail, err)
		return errors.GetErrorWithCode(err)
	}
	return nil
}

func (s *server) SetupActivities(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := struct {
		Hours int `json:"hours"`
	}{}
	_ = decodeRequest(ctx, r, &req)
	if req.Hours == 0 {
		req.Hours = 48
	}
	in2days := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	subC, _ := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	err := subC.SetupActivities(in2days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupActivities(Date:%v). Err:%+v", in2days, err)
	}
	return nil
}

// SetupTags sets up tags for culture preview email and culture email 2 weeks in advance.
func (s *server) SetupTags(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	nextCultureDate := time.Now().Add(time.Hour * 7 * 24)
	for nextCultureDate.Weekday() != time.Monday {
		nextCultureDate = nextCultureDate.Add(24 * time.Hour)
	}
	nextPreviewDate := nextCultureDate
	mailC, err := mail.NewClient(ctx, log, s.serverInfo)
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

// SendPreviewCultureEmail sends the preview email to all subscribers who are not skipped.
func (s *server) SendPreviewCultureEmail(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	cultureDate := time.Now().Add(6 * 24 * time.Hour)
	log.Infof(ctx, "culture date:%s", cultureDate)
	subC := subold.New(ctx)
	subLogs, err := subC.GetForDate(cultureDate)
	if err != nil {
		return errors.Annotate(err, "failed to SendPreviewCultureEmail: failed to subold.GetForDate")
	}
	var nonSkippers []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippers) != 0 {
		if common.IsProd(s.serverInfo.ProjectID) {
			// hard code emails that should be sent email
			nonSkippers = append(nonSkippers, "atish@gigamunchapp.com", "chris@eatgigamunch.com", "enis@eatgigamunch.com", "piyush@eatgigamunch.com", "pkailamanda@gmail.com", "emilywalkerjordan@gmail.com", "mike@eatgigamunch.com", "befutter@gmail.com")
		}
		tag := mail.GetPreviewEmailTag(cultureDate)
		log.Infof(ctx, "culture email tag: %s", tag)
		mailC, err := mail.NewClient(ctx, log, s.serverInfo)
		if err != nil {
			return errors.Annotate(err, "failed to SendPreviewCultureEmail: failed to mail.NewClient")
		}
		err = mailC.AddBatchTags(nonSkippers, []mail.Tag{tag})
		if err != nil {
			return errors.Annotate(err, "failed to SendPreviewCultureEmail: failed to mail.AddBatchTag")
		}
		log.Infof(ctx, "applying tags to: %+v", nonSkippers)
	}
	return nil
}

// SendCultureEmail sends the culture email to all subscribers who are not skipped.
func (s *server) SendCultureEmail(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	cultureDate := time.Now()
	log.Infof(ctx, "culture date: s", cultureDate)
	subC := subold.New(ctx)
	subLogs, err := subC.GetForDate(cultureDate)
	if err != nil {
		return errors.Annotate(err, "failed to SendCultureEmail: failed to subold.GetForDate")
	}
	var nonSkippers []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippers) != 0 {
		if common.IsProd(s.serverInfo.ProjectID) {
			// hard code emails that should be sent email
			nonSkippers = append(nonSkippers, "atish@eatgigamunch.com", "chris@eatgigamunch.com", "enis@eatgigamunch.com", "piyush@eatgigamunch.com", "pkailamanda@gmail.com", "emilywalkerjordan@gmail.com", "mike@eatgigamunch.com", "befutter@gmail.com")
		}
		tag := mail.GetCultureEmailTag(cultureDate)
		log.Infof(ctx, "culture email tag: %s", tag)
		log.Infof(ctx, "applying tags to: %+v", nonSkippers)
		mailC, err := mail.NewClient(ctx, log, s.serverInfo)
		if err != nil {
			return errors.Annotate(err, "failed to SendPreviewCultureEmail: failed to mail.NewClient")
		}
		err = mailC.AddBatchTags(nonSkippers, []mail.Tag{tag})
		if err != nil {
			return errors.Annotate(err, "failed to SendCultureEmail: failed to mail.AddBatchTag")
		}
	}
	return nil
}

// CheckPowerSensors checks all the PowerSensors.
func (s *server) CheckPowerSensors(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	healthC := healthcheck.New(ctx)
	err = healthC.CheckPowerSensors()
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to health.CheckPowerSensors")
	}
	return nil
}

// SendStatsSMS sends the general stats on new and cancel.
func (s *server) SendStatsSMS(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	if !common.IsProd(s.serverInfo.ProjectID) {
		return nil
	}
	subC := subold.New(ctx)
	subs, err := subC.GetHasSubscribed(time.Now())
	if err != nil {
		log.Errorf(ctx, "failed to SendStatsSMS: failed to subold.GetAllSubscribers: %s", err)
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
	for _, s := range subs {
		if s.IsSubscribed {
			totalSubs++
		}
		// week
		if s.SubscriptionDate.After(dateMinus7Days) && s.SubscriptionDate.Before(date) {
			newSubsLastWeek++
		}
		if s.UnSubscribedDate.After(dateMinus7Days) && s.UnSubscribedDate.Before(date) {
			cancelsLastWeek++
			daysWithUs := int(s.UnSubscribedDate.Sub(s.SubscriptionDate) / (time.Hour * 24))
			sumDaysWithUs += daysWithUs
			if daysWithUs > 7*8 {
				cancelsMoreThan8WeekRetention++
			} else if daysWithUs < 7*4 {
				cancelsLessThan4WeekRetention++
			} else {
				cancels4To8WeekRetention++
			}
		}
		if s.SubscriptionDate.Before(dateMinus7Days) {
			totalSubsLastWeek++
		}
		// month
		if s.SubscriptionDate.After(dateMinus30Days) && s.SubscriptionDate.Before(date) {
			newSubs30Days++
		}
		if s.UnSubscribedDate.After(dateMinus30Days) && s.UnSubscribedDate.Before(date) {
			cancelsLast30Days++
		}
		if s.SubscriptionDate.Before(dateMinus30Days) {
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
	for _, number := range numbers {
		err = messageC.SendDeliverySMS(number, msg)
		if err != nil {
			log.Errorf(ctx, "failed to send quantity sms: %+v", err)
		}
	}
	return nil
}

// BackupDatastore creates a back-up datastore in cloud storage.
func (s *server) BackupDatastore(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	// Decode Request
	req := struct {
		Kinds      string `json:"kinds"`
		BucketName string `json:"bucketname"`
	}{}
	_ = decodeRequest(ctx, r, &req)
	log.Infof(ctx, "req: %+v", req)

	projID := s.serverInfo.ProjectID
	// Get Access Token
	accessToken, _, err := appengine.AccessToken(ctx, "https://www.googleapis.com/auth/datastore")
	if err != nil {
		return errInternalError.WithError(err).Annotate("failed to appengine.AccessToken")
	}
	// Create Request
	if req.BucketName == "" {
		req.BucketName = fmt.Sprintf("gs://%s-datastore-backups", projID)
	}
	backupPrefix := req.BucketName
	kinds := strings.Split(req.Kinds, ",")
	type EntityFilter struct {
		Kinds        []string `json:"kinds"`
		NamespaceIDs []string `json:"namespace_ids"`
	}
	entityFilter := EntityFilter{
		Kinds: kinds,
	}
	body := struct {
		ProjectID       string       `json:"project_id"`
		OutputURLPrefix string       `json:"output_url_prefix"`
		EntityFilter    EntityFilter `json:"entity_filter"`
	}{
		ProjectID:       projID,
		OutputURLPrefix: backupPrefix,
		EntityFilter:    entityFilter,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return errInternalError.WithError(err).Annotate("failed to json.Marshal")
	}
	log.Infof(ctx, "backup req: %s", bodyBytes)
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	url := fmt.Sprintf("https://datastore.googleapis.com/v1/projects/%s:export", projID)
	// Make Request
	backupReq, err := http.NewRequest(http.MethodPost, url, bodyBuffer)
	if err != nil {
		return errInternalError.WithError(err).Annotate("failed to http.NewRequest")
	}
	backupReq.Header.Add("Content-Type", "application/json")
	backupReq.Header.Add("Authorization", "Bearer "+accessToken)
	// Results
	result, err := urlfetch.Client(ctx).Do(backupReq)
	if err != nil {
		return errInternalError.WithError(err).Annotate("failed to urlfetch.Client.Do")
	}
	reply, _ := ioutil.ReadAll(result.Body)
	log.Infof(ctx, "Reply: %s", reply)
	return nil
}
