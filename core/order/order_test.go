package order

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/misc/testhelper"
	"github.com/atishpatel/Gigamunch-Backend/types"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

var errResp = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "fake error"}

func getValidMakeOrderReq() MakeOrderReq {
	var exchangeMethod types.ExchangeMethods
	exchangeMethod.SetPickup(true)
	return MakeOrderReq{
		PostID:              1,
		NumServings:         1,
		BTNonce:             "fake-valid-nonce",
		ExchangeMethod:      exchangeMethod,
		GigamuncherAddress:  testhelper.GetGigamuncherAddress(),
		GigamuncherID:       "gigamuncher",
		GigamuncherName:     "Muncher Name",
		GigamuncherPhotoURL: testhelper.PersonPhotoURL,
	}
}

type fakePostClient struct {
	servingsOffered, numServingsOrdered int32
	availableExchangeMethods            types.ExchangeMethods
	closingDateTime                     time.Time
	gigachefDelivery                    post.GigachefDelivery
	AddOrderReq                         *post.AddOrderReq
	getorderinfo, addorder              bool
	GetOrderInfoCalled, AddOrderCalled  bool
}

func newPostClient() *fakePostClient {
	var exchangeMethod types.ExchangeMethods
	exchangeMethod.SetPickup(true)
	return &fakePostClient{
		getorderinfo:             true,
		addorder:                 true,
		servingsOffered:          10,
		numServingsOrdered:       5,
		closingDateTime:          time.Now().Add(time.Minute * 10),
		availableExchangeMethods: exchangeMethod,
		gigachefDelivery: post.GigachefDelivery{
			Price:         10.0,
			Radius:        100,
			TotalDuration: 3600,
		},
	}

}

func (f *fakePostClient) GetOrderInfo(id int64) (*post.OrderInfoResp, error) {
	f.GetOrderInfoCalled = true
	if f.getorderinfo {
		return &post.OrderInfoResp{
			GigachefID:               "chef",
			ItemID:                   10,
			PostTitle:                "title",
			PostPhotoURL:             testhelper.FoodPhotoURL,
			BTSubMerchantID:          "submerch",
			ServingsOffered:          f.servingsOffered,
			NumServingsOrdered:       f.numServingsOrdered,
			ChefPricePerServing:      10,
			PricePerServing:          10,
			TaxPercentage:            7.5,
			AvailableExchangeMethods: f.availableExchangeMethods,
			ClosingDateTime:          f.closingDateTime,
			GigachefDelivery:         f.gigachefDelivery,
			GigachefAddress:          testhelper.GetGigamuncherAddress(),
		}, nil
	}
	return nil, errResp
}

func (f *fakePostClient) AddOrder(ctx context.Context, p *post.AddOrderReq) error {
	f.AddOrderCalled = true
	f.AddOrderReq = p
	if f.addorder {

		return nil
	}
	return errResp
}

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

func TestMakeOrder(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	postC := newPostClient()
	paymentC := newPaymentClient()
	req := getValidMakeOrderReq()
	_, err = makeOrder(ctx, &req, postC, paymentC)
	if err != nil {
		t.Fatal("MakeOrder returned nil with error: ", err)
	}
	if !paymentC.MakeSaleCalled {
		t.Fatal("MakeSale wasn't called.")
	}
	if paymentC.RefundSaleCalled {
		t.Fatal("RefundSale was called.")
	}
	if !postC.GetOrderInfoCalled {
		t.Fatal("GetOrderInfo wasn't called.")
	}
	if !postC.AddOrderCalled {
		t.Fatal("AddOrder wasn't called.")
	}
	tmp := new(Order)
	err = get(ctx, postC.AddOrderReq.OrderID, tmp)
	if err != nil {
		t.Fatal("Order was not created: ", err)
	}
}

func TestMakeOrderThenRefund(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	postC := newPostClient()
	postC.addorder = false
	paymentC := newPaymentClient()
	req := getValidMakeOrderReq()
	_, err = makeOrder(ctx, &req, postC, paymentC)
	if !errResp.Equal(err) {
		t.Fatal("makeOrder: wanted: InternalServerError. received: ", err)
	}
	if !paymentC.MakeSaleCalled {
		t.Fatal("MakeSale wasn't called.")
	}
	if !paymentC.RefundSaleCalled {
		t.Fatal("RefundSale wasn't called.")
	}
	if !postC.GetOrderInfoCalled {
		t.Fatal("GetOrderInfo wasn't called.")
	}
	if !postC.AddOrderCalled {
		t.Fatal("AddOrder wasn't called.")
	}
	tmp := new(Order)
	err = get(ctx, postC.AddOrderReq.OrderID, tmp)
	if err != datastore.ErrNoSuchEntity {
		t.Fatal("Order was not delete: ", err)
	}
}
