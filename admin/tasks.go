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

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"

	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// TODO: UpdateDrip

// SendSMS sends an Customer an SMS from Gigamunch to number.
func (s *server) SendSMS(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	params, err := tasks.ParseSendSMSParam(r)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseSendSMS. Err:%+v", err)
	}
	log.Infof(ctx, "params: %+v", params)
	resp := &pb.ErrorOnlyResp{}
	nameDilm := "{{name}}"
	firstNameDilm := "{{first_name}}"
	emailDilm := "{{email}}"
	userIDDilm := "{{user_id}}"

	messageC := message.New(ctx)
	subC := subold.New(ctx)
	var subs []*subold.SubscriptionSignUp
	var sb *subold.SubscriptionSignUp
	if params.Email != "" {
		sb, err = subC.GetSubscriber(params.Email)
	} else {
		subs, err = subC.GetSubscribersByPhoneNumber(params.Number)
	}
	if err != nil {
		log.Errorf(ctx, "failed to subold.GetSubscriber: %+v", params)
		return errors.Annotate(err, "failed to subold.GetSubscribers. No SMS was sent.")
	}
	if len(subs) > 0 {
		sb = subs[0]
	}
	if sb == nil {
		if params.Number != "" {
			err = messageC.SendDeliverySMS(params.Number, params.Message)
			if err != nil {
				return errors.Annotate(err, "failed to message.SendDeliverySMS to none sub number To:"+params.Number)
			}
			log.Infof(ctx, "message sent for: %+v", params)
			return resp
		}
		log.Errorf(ctx, "failed to GetSubscriber: %+v", params)
		return errors.Annotate(err, "failed to GetSubscriber")
	}
	if sb.PhoneNumber == "" {
		log.Infof(ctx, "no phone number")
		return resp
	}
	name := sb.FirstName
	if name == "" {
		name = sb.Name
	}
	name = strings.Title(name)
	msg := params.Message
	msg = strings.Replace(msg, nameDilm, name, -1)
	msg = strings.Replace(msg, firstNameDilm, sb.FirstName, -1)
	msg = strings.Replace(msg, emailDilm, sb.Email, -1)
	msg = strings.Replace(msg, userIDDilm, sb.ID, -1)
	err = messageC.SendDeliverySMS(sb.PhoneNumber, msg)
	if err != nil {
		return errors.Annotate(err, "failed to message.SendDeliverySMS To:"+sb.PhoneNumber)
	}
	// log
	log.Infof(ctx, "message sent for: %+v", params)
	payload := &logging.MessagePayload{
		Platform: "SMS",
		Body:     msg,
		From:     "Gigamunch",
		To:       sb.PhoneNumber,
	}
	log.SubMessage(sb.ID, sb.Email, payload)

	return resp
}

func (s *server) ProcessActivityTask(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	parms, err := tasks.ParseProcessSubscriptionRequest(r)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseProcessSubscriptionRequest. Err:%+v", err)
	}

	// activityC, _ := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	// err = activityC.Process(parms.Date, parms.SubEmail)
	subC, _ := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	err = subC.ProcessActivity(parms.Date, parms.SubEmail)
	if err != nil {
		log.Errorf(ctx, "failed to activity.Process(Date:%s SubEmail:%s). Err:%+v", parms.Date, parms.SubEmail, err)
		return errors.GetErrorWithCode(err)
	}
	return nil
}

func sendAndLogMessage(ctx context.Context, log *logging.Client, msg string, s *subold.Subscriber) bool {
	messageC := message.New(ctx)
	nameDilm := "{{name}}"
	emailDilm := "{{email}}"
	userIDDilm := "{{user_id}}"
	msg = strings.Replace(msg, nameDilm, s.FullName(), -1)
	msg = strings.Replace(msg, emailDilm, s.Email(), -1)
	msg = strings.Replace(msg, userIDDilm, s.ID, -1)
	number := s.PhoneNumber()
	if number == "" {
		return false
	}
	err := messageC.SendDeliverySMS(number, msg)
	if err != nil {
		log.Errorf(ctx, "failed to message.SendDeliverySMS To(%s): %+v", number, err)
		return false
	}
	// log sms send
	payload := &logging.MessagePayload{
		Platform: "SMS",
		Body:     msg,
		From:     "Gigamunch",
		To:       number,
	}
	log.SubMessage(s.ID, s.Email(), payload)
	return true
}

type hoursReq struct {
	Hours int `json:"hours"`
}

// ProcessUnpaidPreDelivery 3 days before delivery
func (s *server) ProcessUnpaidPreDelivery(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(hoursReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.Hours == 0 {
		req.Hours = 3 * 24
	}

	const smsMessage = "Hey {{name}}, this is Chris from Gigamunch. Just a heads up, you will not receive this week's meal because your account has been suspended due to multiple outstanding charges. Please update your card and settle from declined transactions, we'd love to have you back! Feel free to respond if you have any questions, we'll be happy to clarify. Thank you \n https://eatgigamunch.com/update-payment?email={{email}}"
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to activity.NewClient. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	summaries, err := activityC.GetUnpaidSummaries()
	if err != nil {
		log.Errorf(ctx, "failed to activity.Process. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	// TODO: Try charging card?
	// get subs
	var threePlusUnpaidUserID []string
	for _, v := range summaries {
		switch v.NumUnpaid {
		case "":
		case "0":
		case "1":
		case "2":
			continue
		default:
			// 3=+
			threePlusUnpaidUserID = append(threePlusUnpaidUserID, v.UserID)
		}
	}
	log.Infof(ctx, "number of 3+ unpaid / unpaid subs: %d / %d", len(threePlusUnpaidUserID), len(summaries))
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to sub.NewClient. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	subs, err := subC.GetMulti(threePlusUnpaidUserID)
	if err != nil {
		log.Errorf(ctx, "failed to sub.GetMulti. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	deliveryDay := time.Now().Add(time.Duration(req.Hours) * time.Hour)
	log.Infof(ctx, "Delivery day: %s", deliveryDay)
	for _, s := range subs {
		// check if still subscribed and is for the correct plan day
		if s.Active && s.PlanWeekday == deliveryDay.Weekday().String() {
			nextAct, err := activityC.Get(deliveryDay, s.ID)
			if err != nil {
				log.Errorf(ctx, "failed to process %s: %+v skipping...", s.Email(), err)
				continue
			}
			if !nextAct.Skip {
				// Skip next
				err = activityC.Skip(deliveryDay, s.ID, "Unpaid count")
				if err != nil {
					log.Errorf(ctx, "failed to process %s: failed to skip: %+v skipping...", s.Email(), err)
					continue
				}
				// Send sms
				sendAndLogMessage(ctx, log, smsMessage, s)
			}
		}
	}

	return nil
}

// ProcessUnpaidAutocharge tries to charge unpaid activities.
func (s *server) ProcessUnpaidAutocharge(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to activity.NewClient. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	summaries, err := activityC.GetUnpaidSummaries()
	if err != nil {
		log.Errorf(ctx, "failed to activity.Process. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	// TODO: If older than a year, forgive?
	tasksC := tasks.New(ctx)
	for i := range summaries {
		req := &tasks.ProcessSubscriptionParams{
			UserID:   summaries[i].UserID,
			SubEmail: summaries[i].Email,
			Date:     summaries[i].MaxDateTime(),
		}
		err = tasksC.AddProcessSubscription(time.Now(), req)
		if err != nil {
			log.Errorf(ctx, "failed to AddProcessSubscription: %+v", err)
		}
	}
	return nil
}

// ProcessUnpaidPostDelivery 1 day after delivery
func (s *server) ProcessUnpaidPostDelivery(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(hoursReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.Hours == 0 {
		req.Hours = 24
	}

	const smsMessageTwo = "Hey {{name}}, please don't dine and dash on a local small business. You owe Gigamunch money from declined transactions. Please respond if you have any questions, we'll be happy to clarify. Just tap the link below to settle your outstanding charges. Thank you \n https://eatgigamunch.com/update-payment?email={{email}}"
	const smsMessageThreePlus = "{{name}}, we had to suspend your account for this week because your card was declined 3 times in a row. Please tap the link below to settle your outstanding charges. Feel free to respond if you have any questions, we'll be happy to clarify. Thank you! \n https://eatgigamunch.com/update-payment?email={{email}}"
	const smsMessageDeactive = "Hey {{name}}, please don't dine and dash on a local small business. You owe Gigamunch money from declined transactions. Please respond if you have any questions, we'll be happy to clarify. Just tap the link below to settle your outstanding charges. Thank you \n https://eatgigamunch.com/update-payment?email={{email}}"

	activityC, err := activity.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to activity.NewClient. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	summaries, err := activityC.GetUnpaidSummaries()
	if err != nil {
		log.Errorf(ctx, "failed to activity.Process. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	var userIDs []string
	for _, v := range summaries {
		userIDs = append(userIDs, v.UserID)
	}
	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to sub.NewClient. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}
	subs, err := subC.GetMulti(userIDs)
	if err != nil {
		log.Errorf(ctx, "failed to sub.GetMulti. Err:%+v", err)
		return errors.GetErrorWithCode(err)
	}

	var mailOutstandCharges []string
	var mailUpdateCC []string
	var mailUpdateCCSerious []string

	deliveryDay := time.Now().Add(-1 * time.Duration(req.Hours) * time.Hour)
	log.Infof(ctx, "Delivery day: %s", deliveryDay)

	for i, s := range subs {
		if s.PlanWeekday != deliveryDay.Weekday().String() {
			// wrong day for sub
			continue
		}
		if !s.Active {
			// deactive sub
			mailOutstandCharges = append(mailOutstandCharges, s.Email())
			sendAndLogMessage(ctx, log, smsMessageDeactive, s)
			continue
		} else {
			// active sub
			switch summaries[i].NumUnpaid {
			case "1":
				mailUpdateCC = append(mailUpdateCC, s.Email())
			case "2":
				mailUpdateCC = append(mailUpdateCC, s.Email())
				sendAndLogMessage(ctx, log, smsMessageTwo, s)
			default:
				// 3+
				mailUpdateCCSerious = append(mailUpdateCCSerious, s.Email())
				sendAndLogMessage(ctx, log, smsMessageThreePlus, s)
			}
		}
	}
	log.Infof(ctx, "number of unpaid subs: %d / %d", len(mailOutstandCharges)+len(mailUpdateCC)+len(mailUpdateCCSerious), len(subs))
	// send mail
	log.Infof(ctx, "mailOutstandCharges: %s", mailOutstandCharges)
	log.Infof(ctx, "mailUpdateCC: %s", mailUpdateCC)
	log.Infof(ctx, "mailUpdateCCSerious: %s", mailUpdateCCSerious)
	mailC, err := mail.NewClient(ctx, log, s.serverInfo)
	if err != nil {
		log.Errorf(ctx, "failed to mail.NewClient: %+v", err)
		return errors.GetErrorWithCode(err)
	}
	err = mailC.AddBatchTags(mailOutstandCharges, []mail.Tag{mail.HasOutstandingCharge}, true)
	if err != nil {
		log.Errorf(ctx, "failed to mail.AddBatchTags: %+v", err)
		return errors.GetErrorWithCode(err)
	}
	err = mailC.AddBatchTags(mailUpdateCC, []mail.Tag{mail.UpdateCreditCard}, false)
	if err != nil {
		log.Errorf(ctx, "failed to mail.AddBatchTags: %+v", err)
		return errors.GetErrorWithCode(err)
	}
	err = mailC.AddBatchTags(mailUpdateCCSerious, []mail.Tag{mail.UpdateCreditCardSerious}, false)
	if err != nil {
		log.Errorf(ctx, "failed to mail.AddBatchTags: %+v", err)
		return errors.GetErrorWithCode(err)
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
	nextCultureDateThursday := nextCultureDate.Add(time.Hour * 3 * 24)
	nextPreviewDateThursday := nextCultureDateThursday
	mailC, err := mail.NewClient(ctx, log, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:   "atish@gigamunchapp.com",
		AddTags: []mail.Tag{mail.GetCultureEmailTag(nextCultureDate), mail.GetPreviewEmailTag(nextPreviewDate), mail.GetCultureEmailTag(nextCultureDateThursday), mail.GetPreviewEmailTag(nextPreviewDateThursday)},
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
		err = mailC.AddBatchTags(nonSkippers, []mail.Tag{tag}, false)
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
	log.Infof(ctx, "culture date: %s", cultureDate)
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
		err = mailC.AddBatchTags(nonSkippers, []mail.Tag{tag}, false)
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
		if s.SubscriptionDate.Before(dateMinus7Days) && s.UnSubscribedDate.IsZero() {
			totalSubsLastWeek++
		}
		// month
		if s.SubscriptionDate.After(dateMinus30Days) && s.SubscriptionDate.Before(date) {
			newSubs30Days++
		}
		if s.UnSubscribedDate.After(dateMinus30Days) && s.UnSubscribedDate.Before(date) {
			cancelsLast30Days++
		}
		if s.SubscriptionDate.Before(dateMinus30Days) && s.UnSubscribedDate.IsZero() {
			totalSubs30DaysAgo++
		}
	}
	weeklyChurn := (float32(cancelsLastWeek) / float32(totalSubsLastWeek+cancelsLastWeek)) * 100
	monthlyChurn := (float32(cancelsLast30Days) / float32(totalSubs30DaysAgo+cancelsLast30Days)) * 100
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
	result, err := http.DefaultClient.Do(backupReq)
	if err != nil {
		return errInternalError.WithError(err).Annotate("failed to http.Do")
	}
	reply, _ := ioutil.ReadAll(result.Body)
	log.Infof(ctx, "Reply: %s", reply)
	return nil
}
