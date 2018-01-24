package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	domainURL string
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
			projectID = os.Getenv("PROJECTID")
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
		projectID = os.Getenv("PROJECTID")
	}
	getDomainString()
	http.HandleFunc(tasks.ProcessInquiryURL, handleProcessInquiry)
	http.HandleFunc("/sub-merchant-approved", handleSubMerchantApproved)
	http.HandleFunc("/sub-merchant-declined", handleSubMerchantDeclined)
	http.HandleFunc("/sub-merchant-disbursement-exception", handleDisbursementException)
	// http.HandleFunc("/on-message-sent", handleOnMessageSent)

	http.HandleFunc(tasks.UpdateDripURL, handleUpdateDrip)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)
	http.HandleFunc("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/send-bag-reminder", handleSendBagReminder)
	http.HandleFunc("/send-preview-culture-email", handleSendPreviewCultureEmail)
	http.HandleFunc("/send-culture-email", handleSendCultureEmail)
	http.HandleFunc("/task/process-subscribers", handelProcessSubscribers)
	http.HandleFunc("/task/process-subscription", handelProcessSubscription)
	http.HandleFunc("/task/send-bag-reminder", handleSendBagReminder)
	http.HandleFunc("/task/send-preview-culture-email", handleSendPreviewCultureEmail)
	http.HandleFunc("/task/send-culture-email", handleSendCultureEmail)
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
	register("CreatePromoCode", "createPromoCode", "POST", "cookservice/createPromoCode", "Admin func.")
	register("SetupSubLogs", "SetupSubLogs", "POST", "cookservice/SetupSubLogs", "Setup subscription activty for a date. Admin func. Do this one Chris.")
	register("ProcessSubLog", "ProcessSubLog", "POST", "cookservice/ProcessSubLog", "Admin func.")
	register("CancelSub", "CancelSub", "POST", "cookservice/CancelSub", "Admin func.")
	register("GetSubEmails", "getSubEmails", "GET", "cookservice/getSubEmails", "Admin func.")
	register("SkipSubLog", "skipSubLog", "POST", "cookservice/skipSubLog", "Admin func.")
	register("RefundAndSkipSubLog", "refundAndSkipSubLog", "POST", "cookservice/refundAndSkipSubLog", "Refund and skip a customer for date. Admin func.")
	// register("FreeSubLog", "freeSubLog", "POST", "cookservice/freeSubLog", "Give free meal to a customer for a date. Admin func.")
	register("DiscountSubLog", "DiscountSubLog", "POST", "cookservice/DiscountSubLog", "Give discount to customer. Admin func. ")
	register("ChangeServingsForDate", "ChangeServingsForDate", "POST", "cookservice/ChangeServingsForDate", "Change number of servings for a week. Admin func.")
	register("UpdateMailCustomerFields", "UpdateMailCustomerFields", "POST", "cookservice/UpdateMailCustomerFields", "Updates the custom filds in Drip. Admin func.")
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

func getDomainString() {
	if domainURL == "" {
		domainURL = os.Getenv("DOMAIN_URL")
	}
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
		utils.Criticalf(ctx, "failed to sub.Process(Date:%s SubEmail:%s). \n\nErr:%+v", parms.Date.Format("2006-01-02"), parms.SubEmail, err)
		// TODO schedule for later?
		return
	}
}

func handelProcessSubscribers(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	in6days := time.Now().Add(144 * time.Hour)
	subC := sub.New(ctx)
	err := subC.SetupSubLogs(in6days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupSubLogs(Date:%v). Err:%+v", in6days, err)
		return
	}
}

func handleSendPreviewCultureEmail(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	cultureDate := time.Now().Add(6 * 24 * time.Hour)
	utils.Infof(ctx, "culture date:%s", cultureDate)
	subC := sub.New(ctx)
	subLogs, err := subC.GetForDate(cultureDate)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleSendPreviewCultureEmail: failed to sub.GetForDate: %s", err)
		return
	}
	var nonSkippers []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippers) != 0 {
		if common.IsProd(projectID) {
			// hard code emails that should be sent email
			nonSkippers = append(nonSkippers, "atish@eatgigamunch.com", "chris@eatgigamunch.com", "enis@eatgigamunch.com", "piyush@eatgigamunch.com", "pkailamanda@gmail.com")
		}
		tag := mail.GetPreviewEmailTag(cultureDate)
		mailC := mail.New(ctx)
		err := mailC.AddBatchTags(nonSkippers, []mail.Tag{tag})
		if err != nil {
			utils.Criticalf(ctx, "failed to handleSendPreviewCultureEmail: failed to mail.AddBatchTag: %s", err)
		}
	}
}

func handleSendCultureEmail(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	cultureDate := time.Now()
	subC := sub.New(ctx)
	subLogs, err := subC.GetForDate(cultureDate)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleSendPreviewCultureEmail: failed to sub.GetForDate: %s", err)
		return
	}
	var nonSkippers []string
	for i := range subLogs {
		if !subLogs[i].Skip {
			nonSkippers = append(nonSkippers, subLogs[i].SubEmail)
		}
	}
	if len(nonSkippers) != 0 {
		if common.IsProd(projectID) {
			// hard code emails that should be sent email
			nonSkippers = append(nonSkippers, "atish@eatgigamunch.com", "chris@eatgigamunch.com", "enis@eatgigamunch.com", "piyush@eatgigamunch.com", "pkailamanda@gmail.com")
		}
		tag := mail.GetCultureEmailTag(cultureDate)
		mailC := mail.New(ctx)
		err := mailC.AddBatchTags(nonSkippers, []mail.Tag{tag})
		if err != nil {
			utils.Criticalf(ctx, "failed to handleSendCultureEmail: failed to mail.AddBatchTag: %s", err)
		}
	}
}

func handleSendQuantitySMS(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	if !common.IsProd(projectID) {
		return
	}
	cultureDate := time.Now()
	if cultureDate.Weekday() != time.Monday {
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
	subs, err := subC.GetSubscribers(nonSkippersEmails)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleSendQuantitySMS: failed to sub.GetSubscribers: %s", err)
		return
	}
	twoBags := 0
	twoVegBags := 0
	fourBags := 0
	fourVegBags := 0
	moreThanFourBags := 0
	moreThanFourVegBags := 0
	var specialNames []string
	var listOfMoreThanFourBags []int8
	var listOfMoreThanFourVegBags []int8
	for i, sublog := range nonSkippers {
		veg := false
		if subs[i].VegetarianServings > 0 {
			veg = true
		}
		if sublog.Servings == 2 {
			if veg {
				twoVegBags++
			} else {
				twoBags++
			}
		} else if sublog.Servings == 4 {
			if veg {
				fourVegBags++
			} else {
				fourBags++
			}
		} else {
			if veg {
				moreThanFourVegBags++
				listOfMoreThanFourVegBags = append(listOfMoreThanFourVegBags, sublog.Servings)
			} else {
				moreThanFourBags++
				listOfMoreThanFourBags = append(listOfMoreThanFourBags, sublog.Servings)
			}
		}
		if sublog.Free {
			var name string
			if veg {
				name = fmt.Sprintf("%s %d veg", subs[i].Name, sublog.Servings)
			} else {
				name = fmt.Sprintf("%s %d non-veg", subs[i].Name, sublog.Servings)
			}
			if subs[i].NumGiftDinners > 0 {
				gifter, err := subC.GetSubscriber(subs[i].ReferenceEmail)
				if err == nil {
					name += " gifted from " + gifter.Name
				}
			}
			specialNames = append(specialNames, name)
		}
	}
	totalStandardBags := twoBags + twoVegBags + fourBags + fourVegBags
	msg := `%s culture execution: 
	2 bags: %d 
	2 veg bags: %d
	4 bags: %d 
	4 veg bags: %d

	Total bags: %d + 4+ bags below
	
	4+ bags: %d 
	4+ bags list: %v 
	4+ veg bags: %d
	4+ veg bags list: %v
	
	New Customers: %d 
	%v`
	msg = fmt.Sprintf(msg, cultureDate.Format("Jan 2"), twoBags, twoVegBags, fourBags, fourVegBags, totalStandardBags, len(listOfMoreThanFourBags), listOfMoreThanFourBags, len(listOfMoreThanFourVegBags), listOfMoreThanFourVegBags, len(specialNames), common.CommaSeperate(specialNames))
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
	err = messageC.SendDeliverySMS("6155454989", fmt.Sprintf("Customer Message:\nNumber: %s\nName: %s\nEmail: %s\nBody: %s", from, name, email, body))
	if err != nil {
		utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
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
	mailC := mail.New(ctx)
	// make subscriber if date is same as give reveal date
	sub, err := subC.GetSubscriber(params.Email)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to sub.GetSubscriber: %s", err)
		return
	}
	timeTillGiftReveal := sub.GiftRevealDate.Sub(time.Now())
	utils.Infof(ctx, "time till gift reaveal", timeTillGiftReveal)
	if sub.IsSubscribed && !sub.GiftRevealDate.IsZero() && timeTillGiftReveal < time.Hour*12 {
		mailReq := &mail.UserFields{
			Email:             sub.Email,
			Name:              sub.Name,
			FirstName:         sub.FirstName,
			LastName:          sub.LastName,
			FirstDeliveryDate: sub.FirstBoxDate,
			GifterName:        sub.Reference,
			GifterEmail:       sub.ReferenceEmail,
			AddTags:           []mail.Tag{mail.Subscribed, mail.Customer, mail.Gifted},
		}
		if sub.VegetarianServings > 0 {
			mailReq.AddTags = append(mailReq.AddTags, mail.Vegetarian)
			mailReq.RemoveTags = append(mailReq.RemoveTags, mail.NonVegetarian)
		} else {
			mailReq.AddTags = append(mailReq.AddTags, mail.NonVegetarian)
			mailReq.RemoveTags = append(mailReq.RemoveTags, mail.Vegetarian)
		}
		mailReq.AddTags = append(mailReq.AddTags, mail.GetPreviewEmailTag(sub.FirstBoxDate))
		err = mailC.UpdateUser(mailReq, projectID)
		if err != nil {
			utils.Criticalf(ctx, "Failed to mail.UpdateUser email(%s). Err: %+v", sub.Email, err)
		}
	}
	// add num meals recieved
	activites, err := subC.GetSubscriberActivities(params.Email)
	if err != nil {
		utils.Criticalf(ctx, "failed to handleUpdateDrip: failed to sub.GetForDate: %s", err)
		return
	}
	var numNonSkips int
	for _, activity := range activites {
		if !activity.Skip && activity.Date.Before(time.Now().Add(time.Hour*48)) {
			numNonSkips++
		}
	}
	if numNonSkips >= 1 && numNonSkips <= 3 {
		tag := mail.GetReceivedJourneyTag(numNonSkips)
		utils.Infof(ctx, "Applying Tag(%s) to Email(%s)", tag, params.Email)

		err := mailC.AddTag(params.Email, tag)
		if err != nil {
			logging.Errorf(ctx, "failed to handleUpdateDrip: failed to mail.AddTag: %+v", err)
			w.WriteHeader(500)
			return
		}
	}
	// send chris a message if user reached their set amout of gift meals
	if numNonSkips == sub.NumGiftDinners {
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
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.31623169903713, Longitude: -86.56951904296875}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.34057185894721, Longitude: -86.68075561523438}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.33946565299958, Longitude: -86.73431396484375}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2675285739382, Longitude: -86.75491333007812}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.22544232423855, Longitude: -86.87713623046875}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.140092827322654, Longitude: -86.89773559570312}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.060201412392914, Longitude: -87.05703735351562}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.94243575255426, Longitude: -87.00210571289062}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.870134015336994, Longitude: -87.03506469726562}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.830061559034036, Longitude: -87.04605102539062}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.828948146199636, Longitude: -86.94580078125}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.81669957403484, Longitude: -86.84829711914062}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.821153818963175, Longitude: -86.77001953125}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.84898718690659, Longitude: -86.67938232421875}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.943547570924665, Longitude: -86.649169921875}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.019114512959, Longitude: -86.627197265625}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.049098959065645, Longitude: -86.60247802734375}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.06575205170711, Longitude: -86.55990600585938}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.134547437460064, Longitude: -86.59423828125}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.245380741380465, Longitude: -86.58737182617188}},
			geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.29741818650811, Longitude: -86.55441284179688}},
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
