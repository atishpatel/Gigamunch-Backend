package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/appengine"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
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
	getDomainString()
	http.HandleFunc(tasks.ProcessInquiryURL, handleProcessInquiry)
	http.HandleFunc("/sub-merchant-approved", handleSubMerchantApproved)
	http.HandleFunc("/sub-merchant-declined", handleSubMerchantDeclined)
	http.HandleFunc("/sub-merchant-disbursement-exception", handleDisbursementException)
	http.HandleFunc("/on-message-sent", handleOnMessageSent)
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
		utils.Criticalf(ctx, "Error doing %s: %v", "SubMerchantApproved", err)
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
	transactionIDs, err := paymentC.DisbursementException(signature, payload)
	if err != nil {
		utils.Criticalf(ctx, "Error doing %s: %v", "SubMerchantApproved", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// set inquiries to fulfilled
	inquiryC := inquiry.New(ctx)
	_, err = inquiryC.SetToFulfilledByTransactionIDs(transactionIDs)
	if err != nil {
		utils.Criticalf(ctx, "Error doing %s: %v", "inquiry.SetToFulfilledByTransactionIDs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleOnMessageSent(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	err := req.ParseForm()
	if err != nil {
		utils.Criticalf(ctx, "Error parsing %s request form: %v", "OnMessageSent", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	eventType := req.FormValue("EventType")
	if eventType != "onMessageSent" {
		utils.Warningf(ctx, "invalid event type for handleOnMessageSent: eventType: %+v", eventType)
	}
	channelSid := req.FormValue("ChannelSid")
	body := req.FormValue("Body")
	from := req.FormValue("From")
	if body == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	messageC := message.New(ctx)
	resp, err := messageC.GetChannelInfo(channelSid)
	if err != nil {
		utils.Criticalf(ctx, "errors while to message.GetChannelInfo onMessageSent. err: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userInfo, err := messageC.GetUserInfo(from)
	if err != nil {
		utils.Criticalf(ctx, "errors while to message.GetUserInfo onMessageSent. err: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shouldNotifyCook := resp.CookID != userInfo.ID // isn't cook
	if shouldNotifyCook {
		cookC := cook.New(ctx)
		encodedChannelName := resp.CookID + "%3C%3B%3E" + resp.EaterID
		msg := fmt.Sprintf("%s just sent you a message on Gigamunch:\n\"%s\"\n\n%s/cook/channel/%s", userInfo.Name, body, domainURL, encodedChannelName)
		err = cookC.Notify(resp.CookID, "You just got a message", msg)
		if err != nil {
			utils.Criticalf(ctx, "failed to cook.Notify cookID(%s) in onMessageSent. err: %+v", resp.CookID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
