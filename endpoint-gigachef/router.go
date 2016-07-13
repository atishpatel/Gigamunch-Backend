package gigachef

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/payment"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/post"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

const (
	removeClosedPostsURL = "/remove-closed-posts"
	closePostQueue       = "close-post"
	closePostURL         = "/close-post"
	notifyMuncherQueue   = "notify-muncher"
	notifyMuncherURL     = "/notify-muncher"
	processOrderQueue    = "process-order"
	processOrderURL      = "/process-order"
)

var (
	errParse = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Bad request."}
)

func init() {
	http.HandleFunc(removeClosedPostsURL, handleRemovePosts)
	http.HandleFunc(closePostURL, handleClosePost)
	http.HandleFunc(notifyMuncherURL, handleNotifyMuncher)
	http.HandleFunc(processOrderURL, handleProcessOrder)

	http.HandleFunc("/sub-merchant-approved", handleSubMerchantApproved)
	http.HandleFunc("/sub-merchant-declined", handleSubMerchantDeclined)
	http.HandleFunc("/sub-merchant-disbursement-exception", handleDisbursementException)
}

func handleNotifyMuncher(w http.ResponseWriter, req *http.Request) {
	// TODO notify muncher
	// check if order isn't canceled
	// send notification

}

func handleProcessOrder(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	orderID, err := parseOrderID(req)
	if err != nil {
		utils.Errorf(ctx, "Failed to parse process order form request. Err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderC := order.New(ctx)
	o, err := orderC.Process(orderID)
	if err != nil {
		utils.Errorf(ctx, "Failed to process order. Err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if o.State == order.State.Issues {
		task := createOrderIDTask(ctx, orderID, time.Now().Add(24*time.Hour))
		_, err = taskqueue.Add(ctx, task, processOrderQueue)
		if err != nil {
			utils.Errorf(ctx, "Failed to add to processOrderQueue. Err: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func parseOrderID(req *http.Request) (int64, error) {
	err := req.ParseForm()
	if err != nil {
		return 0, errParse.WithError(err).Wrap("failed to parse from from request")
	}
	orderIDString := req.FormValue("order_id")
	if orderIDString == "" {
		return 0, errParse.Wrapf("Invalid req orderID: %s", orderIDString)
	}
	orderID, err := strconv.ParseInt(orderIDString, 10, 64)
	if err != nil {
		return 0, errParse.WithError(err).Wrapf("failed to parse orderID(%s)", orderIDString)
	}
	return orderID, nil
}

func handleClosePost(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	postID, chefID, err := parsePostIDAndChefID(req)
	if err != nil {
		utils.Errorf(ctx, "Failed to parse close post form request. Err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// remove live post
	postC := post.New(ctx)
	p, err := postC.ClosePost(postID, chefID)
	if err != nil {
		utils.Errorf(ctx, "Failed to close post. Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(p.Orders) > 0 {
		// notify chef
		chefC := gigachef.New(ctx)
		message := fmt.Sprintf("Hey! %d gigamunchers are starving for your %s! Checkout who they are http://www.gigamunchapp.com/posts. :)",
			len(p.Orders),
			p.Title,
		)
		err = chefC.Notify(chefID, "Your customers are starving - Gigamunch", message)
		if err != nil {
			utils.Errorf(ctx, "Failed to notify chef. Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// get process-order queue tasks
		orderIDs := make([]int64, len(p.Orders))
		muncherIDs := make([]string, len(p.Orders))
		for i := range p.Orders {
			orderIDs[i] = p.Orders[i].OrderID
			muncherIDs[i] = p.Orders[i].GigamuncherID
		}
		orderTasks := createOrderTasks(postID, p.Orders)
		// add order queue
		_, err = taskqueue.AddMulti(ctx, orderTasks, processOrderQueue)
		if err != nil {
			utils.Errorf(ctx, "Error adding tasks to process-order queue: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func parsePostIDAndChefID(req *http.Request) (int64, string, error) {
	err := req.ParseForm()
	if err != nil {
		return 0, "", errParse.WithError(err).Wrap("failed to parse from from request")
	}
	postIDString := req.FormValue("post_id")
	chefID := req.FormValue("gigachef_id")
	if postIDString == "" || chefID == "" {
		return 0, "", errParse.Wrapf("Invalid task for close post. postID: %s chefID: %s", postIDString, chefID)
	}
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		return 0, "", errParse.WithError(err).Wrapf("Failed to parse postID(%s). Err: ", postIDString, err)
	}
	return postID, chefID, nil
}

func createOrderTasks(postID int64, orders []post.OrderPost) []*taskqueue.Task {
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	orderTasks := make([]*taskqueue.Task, len(orders))
	for i, order := range orders {
		v := url.Values{}
		v.Set("order_id", strconv.FormatInt(order.OrderID, 10))
		orderTasks[i] = &taskqueue.Task{
			Path:    processOrderURL,
			Payload: []byte(v.Encode()),
			Header:  h,
			Method:  "POST",
			ETA:     order.ExchangeTime.Add(48 * time.Hour),
		}
	}
	return orderTasks
}

func handleRemovePosts(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	postC := post.New(ctx)
	postIDs, chefIDs, err := postC.GetClosedPosts()
	if err != nil {
		utils.Errorf(ctx, "Error removing closed posts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(postIDs) > 0 {
		log.Println("added postID", postIDs)
		tasks := createClosePostTasks(postIDs, chefIDs)
		_, err = taskqueue.AddMulti(ctx, tasks, closePostQueue)
		if err != nil {
			utils.Errorf(ctx, "Error adding tasks to close-post queue: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func createClosePostTasks(postIDs []int64, chefIDs []string) []*taskqueue.Task {
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	tasks := make([]*taskqueue.Task, len(postIDs))
	for i := range postIDs {
		v := url.Values{}
		v.Set("post_id", strconv.FormatInt(postIDs[i], 10))
		v.Set("gigachef_id", chefIDs[i])
		tasks[i] = &taskqueue.Task{
			Path:    closePostURL,
			Payload: []byte(v.Encode()),
			Header:  h,
			Method:  "POST",
		}
	}
	return tasks
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
		utils.Errorf(ctx, "Error parsing %s request form: %v", fnName, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := req.FormValue("bt_payload")
	signature := req.FormValue("bt_signature")
	utils.Infof(ctx, "payload:%#v signature: %s", payload, signature)
	err = fn(signature, payload)
	if err != nil {
		utils.Errorf(ctx, "Error doing %s: %v", "SubMerchantApproved", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDisbursementException(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)

	err := req.ParseForm()
	if err != nil {
		utils.Errorf(ctx, "Error parsing %s request form: %v", "DisbursementException", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := req.FormValue("bt_payload")
	signature := req.FormValue("bt_signature")
	transactionIDs, err := paymentC.DisbursementException(signature, payload)
	if err != nil {
		utils.Errorf(ctx, "Error doing %s: %v", "SubMerchantApproved", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	orderC := order.New(ctx)
	var hadErr bool
	for _, transactionID := range transactionIDs {
		var orderID int64
		orderID, err = orderC.SetStateToPendingByTransactionID(transactionID)
		if err != nil {
			utils.Errorf(ctx, "Failed to set order with transactionID(%s) to pending state: %v", transactionID, err)
			hadErr = true
		} else {
			task := createOrderIDTask(ctx, orderID, time.Now().Add(48*time.Hour))
			_, err = taskqueue.Add(ctx, task, processOrderQueue)
			if err != nil {
				utils.Criticalf(ctx, "Failed to add to processOrderQueue. Err: %v", err)
				hadErr = true
			}
		}
	}
	if hadErr {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createOrderIDTask(ctx context.Context, orderID int64, eta time.Time) *taskqueue.Task {
	h := make(http.Header)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	return &taskqueue.Task{
		Path:    processOrderURL,
		Payload: []byte(v.Encode()),
		Header:  h,
		Method:  "POST",
		ETA:     eta,
	}
}
