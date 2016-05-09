package order

import (
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"

	"google.golang.org/appengine/aetest"
)

var errResp = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "fake error"}

type fakePaymentClient struct {
	makesale, refundsale             bool
	MakeSaleCalled, RefundSaleCalled bool
}

func newPaymentClient() *fakePaymentClient {
	return &fakePaymentClient{
		makesale:   true,
		refundsale: true,
	}
}

func (f *fakePaymentClient) MakeSale(arg1, arg2 string, arg3, arg4 float32) (string, error) {
	f.MakeSaleCalled = true
	if f.makesale {
		return "success", nil
	}
	return "", errResp
}

func (f *fakePaymentClient) RefundSale(string) (string, error) {
	f.RefundSaleCalled = true
	if f.refundsale {
		return "success", nil
	}
	return "", errResp
}

func getCreateReq() *CreateReq {
	return &CreateReq{
		NumServings:          1,
		PricePerServing:      10,
		ChefPricePerServing:  8,
		ExpectedExchangeTime: time.Now(),
	}
}

func TestCreate(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	paymentC := newPaymentClient()
	req := getCreateReq()
	resp, err := create(ctx, req, paymentC)
	if err != nil {
		t.Fatal("Create returned error: ", err)
	}
	if !paymentC.MakeSaleCalled {
		t.Fatal("MakeSale wasn't called.")
	}
	if paymentC.RefundSaleCalled {
		t.Fatal("RefundSale was called.")
	}
	tmp := new(Order)
	err = get(ctx, resp.ID, tmp)
	if err != nil {
		t.Fatal("Order was not created: ", err)
	}
	tmpResp := Resp{
		ID:    resp.ID,
		Order: *tmp,
	}
	// fix nanosecond off error
	resp.CreatedDateTime = resp.CreatedDateTime.Truncate(time.Second)
	tmpResp.CreatedDateTime = tmpResp.CreatedDateTime.Truncate(time.Second)
	resp.ExpectedExchangeDataTime = resp.ExpectedExchangeDataTime.Truncate(time.Second)
	tmpResp.ExpectedExchangeDataTime = tmpResp.ExpectedExchangeDataTime.Truncate(time.Second)
	if !reflect.DeepEqual(tmpResp, *resp) {
		t.Fatalf("Response order does not equal datastore order. \nResp: %#v \nActual: %#v", *resp, tmpResp)
	}
}

func TestCreateAndRefund(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	paymentC := newPaymentClient()
	req := getCreateReq()
	tmpPutIncomplete := putIncomplete
	putIncomplete = func(ctx context.Context, order *Order) (int64, error) {
		return 0, errResp
	}
	_, err = create(ctx, req, paymentC)
	putIncomplete = tmpPutIncomplete
	if !errResp.Equal(err) {
		t.Fatal("makeOrder: wanted: InternalServerError. received: ", err)
	}
	if !paymentC.MakeSaleCalled {
		t.Fatal("MakeSale wasn't called.")
	}
	if !paymentC.RefundSaleCalled {
		t.Fatal("RefundSale wasn't called.")
	}
}

func TestCancel(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	o := new(Order)
	o.GigamuncherID = "muncher"
	orderID, err := putIncomplete(ctx, o)
	if err != nil {
		t.Fatal("failed to put incomplete")
	}
	paymentC := newPaymentClient()
	_, err = cancel(ctx, "muncher", orderID, paymentC)
	if err != nil {
		t.Fatal("cancelOrder returned error: ", err)
	}
	if !paymentC.RefundSaleCalled {
		t.Fatal("RefundSale wasn't called.")
	}

	tmp := new(Order)
	err = get(ctx, orderID, tmp)
	if err != nil {
		t.Fatal("Error getting order: ", err)
	}
	if tmp.State != State.Refunded {
		t.Fatal("Order state was not 'refunded': ", tmp)
	}
}
