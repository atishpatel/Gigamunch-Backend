package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	projectID string
)

type coder interface {
	GetCode() int32
}

func handleResp(ctx context.Context, fnName string, resp coder) {
	code := resp.GetCode()
	if code == errors.CodeInvalidParameter {
		utils.Warningf(ctx, "%s invalid parameter: %+v", fnName, resp)
		return
	} else if code != 0 {
		if projectID == "" {
			projectID = os.Getenv("PROJECT_ID")
		}
		utils.Criticalf(ctx, "%s COOKAPI: %s err: %+v", projectID, fnName, resp)
	}
}

type validatableTokenReq interface {
	gigatoken() string
	valid() error
}

func validateRequestAndGetUser(ctx context.Context, req validatableTokenReq) (*types.User, error) {
	err := req.valid()
	if err != nil {
		return nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}.Wrap("failed to validate request")
	}
	user, err := auth.GetUserFromToken(ctx, req.gigatoken())
	if err != nil && errors.GetErrorWithCode(err).Code == errors.CodeSignOut {
		utils.Criticalf(ctx, "A signout code was issued to. err: %+v", err)
	}
	return user, err
}

// Service is the REST API Endpoint exposed to Gigamunchers
type Service struct{}

func main() {
	if projectID == "" {
		projectID = os.Getenv("PROJECT_ID")
	}
	http.HandleFunc(tasks.ProcessInquiryURL, handleProcessInquiry)
	http.HandleFunc("/sub-merchant-approved", handleSubMerchantApproved)
	http.HandleFunc("/sub-merchant-declined", handleSubMerchantDeclined)
	http.HandleFunc("/sub-merchant-disbursement-exception", handleDisbursementException)
	// http.HandleFunc("/on-message-sent", handleOnMessageSent)

	http.HandleFunc(tasks.UpdateDripURL, handleUpdateDrip)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)
	http.HandleFunc("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/send-bag-reminder", handleSendBagReminder)
	http.HandleFunc("/task/process-subscribers", handelProcessSubscribers)
	http.HandleFunc("/task/process-subscription", handelProcessSubscription)
	http.HandleFunc("/task/send-bag-reminder", handleSendBagReminder)
	http.HandleFunc("/task/send-quantity-sms", handleSendQuantitySMS)
	http.HandleFunc("/send-quantity-sms", handleSendQuantitySMS)
	http.HandleFunc("/webhook/twilio-sms", handleTwilioSMS)
	http.HandleFunc("/testbra", testbra)
	api, err := endpoints.RegisterService(&Service{}, "cookservice", "v1", "An endpoint service for cooks.", true)
	if err != nil {
		log.Fatalf("Failed to register service: %#v", err)
	}

	register := func(orig, name, method, path, desc string) {
		m := api.MethodByName(orig)
		if m == nil {
			log.Fatalf("Missing method %s", orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
	}
	// Page stuff
	register("FinishOnboarding", "finishOnboarding", "POST", "cookservice/finishOnboarding", "Updates cook and submerchant")
	// Refresh stuff
	register("RefreshToken", "refreshToken", "POST", "cookservice/refreshToken", "Refresh a token.")
	// Cook stuff
	register("GetCook", "getCook", "GET", "cookservice/getCook", "Get the cook info.")
	register("UpdateCook", "updateCook", "POST", "cookservice/updateCook", "Update cook information.")
	register("SchedulePhoneCall", "schedulePhoneCall", "POST", "cookservice/schedulePhoneCall", "Schedule a phone call.")
	// Submerchant stuff
	register("UpdateSubMerchant", "updateSubMerchant", "POST", "gigachefservice/updateSubMerchant", "Update or create sub-merchant.")
	register("GetSubMerchant", "getSubMerchant", "GET", "gigachefservice/getSubMerchant", "Get the sub merchant info.")
	// Item stuff
	register("SaveItem", "saveItem", "POST", "cookservice/saveItem", "Save an item.")
	register("GetItem", "getItem", "GET", "cookservice/getItem", "Get an item.")
	register("ActivateItem", "activateItem", "POST", "cookservice/activateItem", "Activate an item.")
	register("DeactivateItem", "deactivateItem", "POST", "cookservice/deactivateItem", "Deactivate an item.")
	// Menu stuff
	register("GetMenus", "getMenus", "GET", "cookservice/getMenus", "Gets the menus for a cook.")
	register("SaveMenu", "saveMenu", "POST", "cookservice/saveMenu", "Save a menu.")
	// Inquiry stuffffffffff
	register("GetMessageToken", "getMessageToken", "GET", "cookservice/getMessageToken", "Gets the a token for messaging.")
	register("GetInquiries", "getInquiries", "GET", "cookservice/getInquiries", "GetInquiries gets a cook's inquiries.")
	register("GetInquiry", "getInquiry", "GET", "cookservice/getInquiry", "GetInquiry gets a cook's inquiry.")
	register("AcceptInquiry", "acceptInquiry", "POST", "cookservice/acceptInquiry", "AcceptInquiry accepts an inquiry for a cook.")
	register("DeclineInquiry", "declineInquiry", "POST", "cookservice/declineInquiry", "DeclineInquiry declines an inquiry for a cook.")
	// Admin stuffffffffffffffff
	register("AddToProcessInquiryQueue", "addToProcessInquiryQueue", "POST", "cookservice/addToProcessInquiryQueue", "Admin func.")
	register("CreateFakeGigatoken", "createFakeGigatoken", "POST", "cookservice/createFakeGigatoken", "Admin func.")
	register("CreateFakeSubmerchant", "createFakeSubmerchant", "POST", "cookservice/createFakeSubmerchant", "Admin func.")
	register("SendSMS", "sendSMS", "POST", "cookservice/sendSMS", "Admin func.")
	register("SendCustomerSMS", "sendCustomerSMS", "POST", "cookservice/sendCustomerSMS", "Admin func.")
	register("CreatePromoCode", "createPromoCode", "POST", "cookservice/createPromoCode", "Admin func.")
	register("SetupSubLogs", "SetupSubLogs", "POST", "cookservice/SetupSubLogs", "Setup subscription activty for a date. Admin func. Do this one Chris.")
	register("ProcessSubLog", "ProcessSubLog", "POST", "cookservice/ProcessSubLog", "Admin func.")
	register("CancelSub", "CancelSub", "POST", "cookservice/CancelSub", "Admin func.")
	register("GetSubEmails", "getSubEmails", "POST", "cookservice/getSubEmails", "Admin func.")
	register("SkipSubLog", "skipSubLog", "POST", "cookservice/skipSubLog", "Admin func.")
	register("RefundAndSkipSubLog", "refundAndSkipSubLog", "POST", "cookservice/refundAndSkipSubLog", "Refund and skip a customer for date. Admin func.")
	register("GetGeneralStats", "GetGeneralStats", "POST", "cookservice/GetGeneralStats", "Returns general stats. Admin func.")
	// register("FreeSubLog", "freeSubLog", "POST", "cookservice/freeSubLog", "Give free meal to a customer for a date. Admin func.")
	register("DiscountSubLog", "DiscountSubLog", "POST", "cookservice/DiscountSubLog", "Give discount to customer. Admin func. ")
	register("ChangeServingsForDate", "ChangeServingsForDate", "POST", "cookservice/ChangeServingsForDate", "Change number of servings for a week. Admin func.")
	register("UpdatePaymentMethodToken", "UpdatePaymentMethodToken", "POST", "cookservice/UpdatePaymentMethodToken", "Updates the payment method token. Admin func.")
	register("ChangeServingsPermanently", "ChangeServingsPermanently", "POST", "cookservice/ChangeServingsPermanently", "Change number of servings permanently. Admin func.")
	register("GetSubLogs", "getSubLogs", "POST", "cookservice/getSubLogs", "Get all subscription activty. Admin func.")
	register("GetSubLogsForDate", "getSubLogsForDate", "POST", "cookservice/getSubLogsForDate", "Get subscription activty for a date. Admin func.")
	register("AddToProcessSubscriptionQueue", "addToProcessSubscriptionQueue", "POST", "cookservice/addToProcessSubscriptionQueue", "Admin func.")
	// register("SendWelcomeEmail", "SendWelcomeEmail", "POST", "cookservice/SendWelcomeEmail", "Admin func. Sends welcome email.")
	// register("SendIntroEmail", "SendIntroEmail", "POST", "cookservice/SendIntroEmail", "Admin func. Sends email for people who just left email.")
	endpoints.HandleHTTP()
	appengine.Main()
}

func handleProcessInquiry(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	inquiryID, err := tasks.ParseInquiryID(req)
	if err != nil {
		utils.Criticalf(ctx, "Failed to parse process inquiry request. Err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inquiryC := inquiry.New(ctx)
	err = inquiryC.Process(inquiryID)
	if err != nil {
		utils.Criticalf(ctx, "Failed to process inquiry(%d). Err: %v", inquiryID, err)
		taskC := tasks.New(ctx)
		err = taskC.AddProcessInquiry(inquiryID, time.Now().Add(1*time.Hour))
		if err != nil {
			utils.Criticalf(ctx, "Failed to add inquiry(%d) in processInquiry queue. Err: %v", inquiryID, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func handleSubMerchantApproved(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)
	btNotification(ctx, w, req, "SubMerchantApproved", paymentC.SubMerchantApproved)
}

func handleSubMerchantDeclined(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)
	btNotification(ctx, w, req, "SubMerchantDeclined", paymentC.SubMerchantDeclined)
}

func btNotification(ctx context.Context, w http.ResponseWriter, req *http.Request, fnName string, fn func(string, string) error) {
	err := req.ParseForm()
	if err != nil {
		utils.Criticalf(ctx, "Error parsing %s request form: %v", fnName, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := req.FormValue("bt_payload")
	signature := req.FormValue("bt_signature")
	utils.Infof(ctx, "payload:%#v signature: %s", payload, signature)
	err = fn(signature, payload)
	if err != nil {
		utils.Criticalf(ctx, "Error doing %s: %v", fnName, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDisbursementException(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	err := req.ParseForm()
	if err != nil {
		utils.Criticalf(ctx, "Error parsing %s request form: %v", "DisbursementException", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	paymentC := payment.New(ctx)
	payload := req.FormValue("bt_payload")
	signature := req.FormValue("bt_signature")
	_, err = paymentC.DisbursementException(signature, payload)
	if err != nil {
		utils.Criticalf(ctx, "Error doing %s: %v", "DisbursementException", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// set inquiries to fulfilled
	// TODO Braintree auto release when a submerchant updates their banking info?
	// inquiryC := inquiry.New(ctx)
	// _, err = inquiryC.SetToFulfilledByTransactionIDs(transactionIDs)
	// if err != nil {
	// 	utils.Criticalf(ctx, "Error doing %s: %v", "inquiry.SetToFulfilledByTransactionIDs", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
}

// func handleOnMessageSent(w http.ResponseWriter, req *http.Request) {
// 	ctx := appengine.NewContext(req)
// 	err := req.ParseForm()
// 	if err != nil {
// 		utils.Criticalf(ctx, "Error parsing %s request form: %v", "OnMessageSent", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	eventType := req.FormValue("EventType")
// 	if eventType != "onMessageSent" {
// 		utils.Warningf(ctx, "invalid event type for handleOnMessageSent: eventType: %+v", eventType)
// 	}
// 	channelSid := req.FormValue("ChannelSid")
// 	body := req.FormValue("Body")
// 	from := req.FormValue("From")
// 	if body == "" {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// 	messageC := message.New(ctx)
// 	resp, err := messageC.GetChannelInfo(channelSid)
// 	if err != nil {
// 		utils.Criticalf(ctx, "errors while to message.GetChannelInfo onMessageSent. err: %+v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	userInfo, err := messageC.GetUserInfo(from)
// 	if err != nil {
// 		utils.Criticalf(ctx, "errors while to message.GetUserInfo onMessageSent. err: %+v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	shouldNotifyCook := resp.CookID != userInfo.ID // isn't cook
// 	gigamessage := ""
// 	if shouldNotifyCook {
// 		gigamessage = fmt.Sprintf("From: %s\nTo: %s\n Message:%s", userInfo.Name, resp.EaterName, body)
// 	} else {
// 		gigamessage = fmt.Sprintf("From: %s\nTo: %s\n Message:%s", userInfo.Name, resp.CookName, body)
// 	}
// 	err = messageC.SendSMS("6153975516", gigamessage)
// 	if err != nil {
// 		utils.Criticalf(ctx, "failed to notify enis about message on app. err: %s", err)
// 	}
// 	if shouldNotifyCook {
// 		cookC := cook.New(ctx)
// 		encodedChannelName := resp.CookID + "%3C%3B%3E" + resp.EaterID
// 		msg := fmt.Sprintf("%s just sent you a message on Gigamunch:\n\"%s\"\n\n%s/cook/channel/%s", userInfo.Name, body, domainURL, encodedChannelName)
// 		err = cookC.Notify(resp.CookID, "You just got a message", msg)
// 		if err != nil {
// 			utils.Criticalf(ctx, "failed to cook.Notify cookID(%s) in onMessageSent. err: %+v", resp.CookID, err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

func handelProcessSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	parms, err := tasks.ParseProcessSubscriptionRequest(req)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseProcessSubscriptionRequest. Err:%+v", err)
		return
	}
	subC := sub.New(ctx)
	err = subC.Process(parms.Date, parms.SubEmail)
	if err != nil {
		utils.Errorf(ctx, "failed to sub.Process(Date:%s SubEmail:%s). \n\nErr:%+v", parms.Date.Format("2006-01-02"), parms.SubEmail, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handelProcessSubscribers(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	in7days := time.Now().Add(24 * 7 * time.Hour)
	in4days := time.Now().Add(24 * 4 * time.Hour)
	subC := sub.New(ctx)
	err := subC.SetupSubLogs(in7days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupSubLogs(Date:%v). Err:%+v", in7days, err)
		return
	}
	err = subC.SetupSubLogs(in4days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupSubLogs(Date:%v). Err:%+v", in4days, err)
		return
	}
}

func handleSendQuantitySMS(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	if !common.IsProd(projectID) {
		return
	}
	cultureDate := time.Now()
	for cultureDate.Weekday() != time.Monday {
		cultureDate = cultureDate.Add(24 * time.Hour)
	}
	subC := sub.New(ctx)
	subLogs, err := subC.GetForDate(cultureDate)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleSendQuantitySMS: failed to sub.GetForDate: %s", err)
		return
	}
	var nonSkippers []*sub.SubscriptionLog
	var nonSkippersEmails []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i])
			nonSkippersEmails = append(nonSkippersEmails, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippersEmails) == 0 {
		return
	}

	twoBags := 0
	twoVegBags := 0
	fourBags := 0
	fourVegBags := 0
	specialTwoBags := 0
	specialTwoVegBags := 0
	specialFourBags := 0
	specialFourVegBags := 0
	var listOfMoreThanFourBags []int8
	var listOfMoreThanFourVegBags []int8
	for _, sublog := range nonSkippers {
		vegServings := sublog.VegServings
		nonVegServings := sublog.Servings

		switch vegServings {
		case 0:
			// break
		case 2:
			if sublog.Free {
				specialTwoVegBags++
			} else {
				twoVegBags++
			}
		case 4:
			if sublog.Free {
				specialFourVegBags++
			} else {
				fourVegBags++
			}
		default:
			listOfMoreThanFourVegBags = append(listOfMoreThanFourVegBags, vegServings)
		}

		switch nonVegServings {
		case 0:
			// break
		case 2:
			if sublog.Free {
				specialTwoBags++
			} else {
				twoBags++
			}
		case 4:
			if sublog.Free {
				specialFourBags++
			} else {
				fourBags++
			}
		default:
			listOfMoreThanFourBags = append(listOfMoreThanFourBags, nonVegServings)
		}
	}
	for _, special := range listOfMoreThanFourBags {
		for special >= 4 {
			fourBags++
			special -= 4
		}
		for special >= 2 {
			twoBags++
			special -= 2
		}
	}
	for _, special := range listOfMoreThanFourVegBags {
		for special >= 4 {
			fourVegBags++
			special -= 4
		}
		for special >= 2 {
			twoVegBags++
			special -= 2
		}
	}
	totalStandardBags := twoBags + twoVegBags + fourBags + fourVegBags + specialTwoBags + specialFourBags + specialTwoVegBags + specialFourVegBags
	msg := `%s culture execution:
	2 bags: %d
	2 special bags: %d
	4 bags: %d
	4 special bags: %d

	2 veg bags: %d
	2 special veg bags: %d
	4 veg bags: %d
	4 special veg bags: %d

	Total bags: %d

	Accounted for above:
	4+ bags list: %v
	4+ veg bags list: %v`
	msg = fmt.Sprintf(msg, cultureDate.Format("Jan 2"), twoBags, specialTwoBags, fourBags, specialFourBags, twoVegBags, specialTwoVegBags, fourVegBags, specialFourVegBags, totalStandardBags, listOfMoreThanFourBags, listOfMoreThanFourVegBags)
	messageC := message.New(ctx)
	numbers := []string{"9316445311", "6155454989", "6153975516", "9316446755", "6154913694"}
	for _, number := range numbers {
		err = messageC.SendDeliverySMS(number, msg)
		if err != nil {
			logging.Errorf(ctx, "failed to send quantity sms: %+v", err)
		}
	}
}

func handleSendBagReminder(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	tomorrow := time.Now().Add(10 * time.Hour) // cron job runs at 8 PM
	subC := sub.New(ctx)
	subLogs, err := subC.GetForDate(tomorrow)
	if err != nil {
		utils.Criticalf(ctx, "failed to SendBagReminder: failed to sub.GetForDate: %s", err)
		return
	}
	var nonSkippers []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippers) != 0 {
		subs, err := subC.GetSubscribers(nonSkippers)
		if err != nil {
			utils.Criticalf(ctx, "failed to SendBagReminder: failed to sub.GetSubscribers: %s", err)
			return
		}
		messageC := message.New(ctx)
		for _, sub := range subs {
			if sub.PhoneNumber != "" && sub.BagReminderSMS {
				err := messageC.SendBagSMS(sub.PhoneNumber, fmt.Sprintf("Hey %s! Friendly reminder to leave your Gigamunch bag out tonight or tomorrow morning. Thank you! ^_^", sub.GetName()))
				if err != nil {
					utils.Criticalf(ctx, "error in SendBagReminder: failed to message.SendSMS to %s: %s", sub.PhoneNumber, err)
					continue
				}
				utils.Infof(ctx, "notifed %s(%s)", sub.Name, sub.PhoneNumber)
			}
		}
	}

}

func handleTwilioSMS(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	err := req.ParseForm()
	utils.Infof(ctx, "req body: %s err: %s", req.Form, err)
	from := req.FormValue("From")
	l := len(from)
	from = from[:l-7] + "-" + from[l-7:l-4] + "-" + from[l-4:]
	body := req.FormValue("Body")
	var name, email string
	// TODO: auto get name and email
	messageC := message.New(ctx)
	if from == "+1615-545-4989" {
		splitBody := strings.Split(body, "::")
		if len(splitBody) < 2 {
			return
		}
		err = messageC.SendDeliverySMS(splitBody[0], splitBody[1])
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to sub. Err: %+v", err)
		}
		err = messageC.SendDeliverySMS("6155454989", fmt.Sprintf("Message successfuly send to %s.", splitBody[0]))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
	} else {
		err = messageC.SendDeliverySMS("6155454989", fmt.Sprintf("Customer Message:\nNumber: %s\nName: %s\nEmail: %s\nBody: %s", from, name, email, body))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
	}
}

func handleUpdateDrip(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	params, err := tasks.ParseUpdateDripRequest(req)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to tasks.ParseUpdateDripRequest: %+v", err)
		return
	}
	logging.Infof(ctx, "Params: %+v", params)

	subC := sub.New(ctx)
	log, serverInfo, _, err := setupLoggingAndServerInfo(ctx, "/cookapi/UpdateDrip")
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to setupLoggingAndServerInfo: %s", err)
		return
	}

	// make subscriber if date is same as give reveal date
	sub, err := subC.GetSubscriber(params.Email)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to sub.GetSubscriber: %s", err)
		return
	}
	if !sub.IsSubscribed {
		utils.Infof(ctx, "%s is not a subscriber", params.Email)
		return
	}
	timeTillGiftReveal := sub.GiftRevealDate.Sub(time.Now())
	if !sub.GiftRevealDate.IsZero() && timeTillGiftReveal > time.Hour*12 {
		utils.Infof(ctx, "%s is not ready for gift reveal", params.Email)
		return
	}

	// add num meals recieved
	activites, err := subC.GetSubscriberActivities(params.Email)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to sub.GetSubscriberActivities: %s", err)
		return
	}
	var numNonSkips int
	for _, activity := range activites {
		if !activity.Skip && activity.Date.Before(time.Now().Add(time.Hour*48)) {
			numNonSkips++
		}
	}

	addTags := []mail.Tag{mail.Subscribed, mail.Subscriber}
	// add gift tag
	if !sub.GiftRevealDate.IsZero() {
		addTags = append(addTags, mail.Gifted)
	}
	// add journey tags
	if numNonSkips >= 1 {
		addTags = append(addTags, mail.GetReceivedJourneyTag(1))
	}
	if numNonSkips >= 2 {
		addTags = append(addTags, mail.GetReceivedJourneyTag(2))
	}
	if numNonSkips >= 3 {
		addTags = append(addTags, mail.GetReceivedJourneyTag(3))
	}
	if numNonSkips >= 5 {
		addTags = append(addTags, mail.GetReceivedJourneyTag(5))
	}

	// Update Drip
	mailReq := &mail.UserFields{
		Email:             sub.Email,
		FirstName:         sub.FirstName,
		LastName:          sub.LastName,
		FirstDeliveryDate: sub.FirstBoxDate,
		GifterName:        sub.Reference,
		GifterEmail:       sub.ReferenceEmail,
		AddTags:           addTags,
		VegServings:       sub.VegetarianServings,
		NonVegServings:    sub.Servings,
	}
	mailC, err := mail.NewClient(ctx, log, serverInfo)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to mail.NewClient: %s", err)
		return
	}
	mailReq.AddTags = append(mailReq.AddTags, mail.GetPreviewEmailTag(sub.FirstBoxDate))
	err = mailC.SubActivated(mailReq)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to mail.SubActivated email(%s). Err: %+v", sub.Email, err)
	}
	// send chris a message if user reached their set amout of gift meals
	if !sub.GiftRevealDate.IsZero() && numNonSkips == sub.NumGiftDinners {
		messageC := message.New(ctx)
		err = messageC.SendAdminSMS("6155454989", fmt.Sprintf("Person is done with their gifted meals \nName: %s\nEmail: %s", sub.Name, sub.Email))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
	}
}

func testbra(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	fence := &geofence.Geofence{
		ID:   "Nashville",
		Type: geofence.ServiceZone,
		Name: "Nashville",
		Points: []geofence.Point{
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3195513, Longitude: -86.5475464}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3347628, Longitude: -86.5248873}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3521861, Longitude: -86.5420532}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3438904, Longitude: -86.7253876}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2719574, Longitude: -86.7576599}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2459345, Longitude: -86.8139649}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2304274, Longitude: -86.8805695}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1971752, Longitude: -86.9011731}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1417563, Longitude: -86.8956756}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.079072, Longitude: -87.0570373}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.9518857, Longitude: -87.0206452}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8253964, Longitude: -87.0253076}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8233808, Longitude: -86.8421173}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8462043, Longitude: -86.6670227}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.9526416, Longitude: -86.6629813}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.0080895, Longitude: -86.6210983}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.0835115, Longitude: -86.5626526}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1323292, Longitude: -86.5956116}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1927545, Longitude: -86.5647125}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2614385, Longitude: -86.5805054}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3195513, Longitude: -86.5475464}},
		},
	}
	key := datastore.NewKey(ctx, "Geofence", "Nashville", 0, nil)
	_, err := datastore.Put(ctx, key, fence)
	if err != nil {
		w.Write([]byte(fmt.Sprintln("fail: ", err)))
		return
	}
	key = datastore.NewKey(ctx, "Geofence", "", common.Nashville.ID(), nil)
	_, err = datastore.Put(ctx, key, fence)
	if err != nil {
		w.Write([]byte(fmt.Sprintln("fail: ", err)))
		return
	}
	w.Write([]byte("success"))
	w.Write([]byte(projectID))
}
