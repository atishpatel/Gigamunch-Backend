package post

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/misc/testhelper"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"

	"google.golang.org/appengine/aetest"
)

var errResp = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "fake error"}

type fakeOrderClient struct {
	resp                                        *order.Resp
	postID                                      int64
	create, cancel, getPostID                   bool
	CreateCalled, CancelCalled, GetPostIDCalled bool
}

func (c *fakeOrderClient) Create(ctx context.Context, req *order.CreateReq) (*order.Resp, error) {
	c.CreateCalled = true
	if c.create {
		return c.resp, nil
	}
	return nil, errResp
}

func (c *fakeOrderClient) Cancel(ctx context.Context, userID string, orderID int64) (*order.Resp, error) {
	c.CancelCalled = true
	if c.cancel {
		return c.resp, nil
	}
	return nil, errResp
}

func (c *fakeOrderClient) GetPostID(orderID int64) (int64, error) {
	c.GetPostIDCalled = true
	if c.getPostID {
		return c.postID, nil
	}
	return 0, errResp
}

func newOrderClient(postID int64) *fakeOrderClient {
	return &fakeOrderClient{
		resp:      &order.Resp{ID: -1},
		postID:    postID,
		create:    true,
		cancel:    true,
		getPostID: true,
	}
}

func getValidMakeOrderReq(ctx context.Context, t *testing.T) *MakeOrderReq {
	var exchangeMethod types.ExchangeMethods
	exchangeMethod.SetPickup(true)
	p := &Post{
		ClosingDateTime:          time.Now().Add(time.Hour),
		ServingsOffered:          10,
		AvailableExchangeMethods: exchangeMethod,
	}
	postID, err := putIncomplete(ctx, p)
	if err != nil {
		t.Fatal("put incomplete failed")
	}
	return &MakeOrderReq{
		PostID:             postID,
		NumServings:        1,
		PaymentNonce:       "valid",
		ExchangeMethod:     exchangeMethod,
		GigamuncherAddress: testhelper.GetGigamuncherAddress(),
		GigamuncherID:      "muncher",
	}
}

func TestMakeOrder(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	req := getValidMakeOrderReq(ctx, t)
	orderC := newOrderClient(req.PostID)
	// test
	_, err = makeOrder(ctx, req, orderC)
	if err != nil {
		t.Fatal("MakeOrder returned error: ", err)
	}
	if !orderC.CreateCalled {
		t.Fatal("orderC.Create wasn't called.")
	}
	if orderC.CancelCalled {
		t.Fatal("orderC.Cancel was called.")
	}
	tmp := new(Post)
	err = get(ctx, req.PostID, tmp)
	if err != nil {
		t.Fatal("error getting post: ", err)
	}
	if len(tmp.Orders) == 0 {
		t.Fatal("Order wasn't added to the post.")
	}
}

func TestCancelOrder(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	p := &Post{
		ClosingDateTime: time.Now().Add(time.Hour),
		ServingsOffered: 10,
		Orders: []OrderPost{
			OrderPost{
				OrderID:       10,
				GigamuncherID: "muncher",
			},
			OrderPost{
				OrderID: 20,
			},
		},
	}

	postID, err := putIncomplete(ctx, p)
	orderC := newOrderClient(postID)
	// test
	_, err = cancelOrder(ctx, "muncher", 10, orderC)
	if err != nil {
		t.Fatal("cancelOrder returned error: ", err)
	}
	if !orderC.CancelCalled {
		t.Fatal("Cancel wasn't called.")
	}
	tmpP := new(Post)
	err = get(ctx, postID, tmpP)
	if err != nil {
		t.Fatal("get post failed ")
	}
	for i := range tmpP.Orders {
		if tmpP.Orders[i].OrderID == 10 {
			t.Fatal("OrderID 10 was still in post after cancel.")
		}
	}
}

func TestCancelOrderErrors(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	p := &Post{
		BaseItem: types.BaseItem{
			GigachefID: "chef",
		},
		ClosingDateTime: time.Now().Add(time.Hour),
		ServingsOffered: 10,
	}

	noOrdersPostID, err := putIncomplete(ctx, p)
	if err != nil {
		t.Fatal("failed to put post")
	}
	p.Orders = []OrderPost{
		OrderPost{
			OrderID:       10,
			GigamuncherID: "muncher",
		},
	}
	postID, err := putIncomplete(ctx, p)
	if err != nil {
		t.Fatal("failed to put post")
	}

	p.ClosingDateTime = time.Now()
	closedPostID, err := putIncomplete(ctx, p)
	if err != nil {
		t.Fatal("failed to put post")
	}
	cases := []struct {
		description string
		userID      string
		orderC      orderClient
		want        errors.ErrorWithCode
	}{
		{
			description: "Order not found",
			userID:      "muncher",
			orderC:      newOrderClient(noOrdersPostID),
			want:        errInvalidParameter,
		},
		{
			description: "Invalid userID",
			userID:      "muc",
			orderC:      newOrderClient(postID),
			want:        errUnauthorized,
		},
		{
			description: "Invalid postID",
			userID:      "muncher",
			orderC:      newOrderClient(0),
			want:        errDatastore,
		},
		{
			description: "Post is closed",
			userID:      "muncher",
			orderC:      newOrderClient(closedPostID),
			want:        errPostIsClosed,
		},
	}
	// test
	for _, test := range cases {
		_, err := cancelOrder(ctx, test.userID, 10, test.orderC)
		if !test.want.Equal(err) {
			t.Errorf("%s test failed | Wanted: %s | Recieved: %s", test.description, test.want, err)
		}
	}
}
