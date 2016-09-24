package payment

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"

	"google.golang.org/appengine/aetest"
)

const (
	validNonce = "fake-valid-nonce"
)

var subMerchantID string

func makeSale(t *testing.T, c Client, amount, serviceFee float32) string {
	s, err := c.MakeSale(subMerchantID, validNonce, amount, serviceFee)
	if err != nil {
		t.Error("failed to make sale: ", err)
	}
	return s
}

func TestMakeSaleThenRelease(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	c := New(ctx)
	transactionID := makeSale(t, c, 2, .1)
	_, err = c.bt.Transaction().Settle(transactionID)
	if err != nil {
		t.Error("failed to settle sale: ", err)
	}
	_, err = c.ReleaseSale(transactionID)
	if err != nil {
		t.Error("failed to release sale: ", err)
	}
}

func TestMakeSaleThenRefund(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// setup
	// test void
	c := New(ctx)
	transactionID := makeSale(t, c, 1, .1)
	_, err = c.RefundSale(transactionID)
	if err != nil {
		t.Error("failed to void sale: ", err)
	}

	// test refund
	transactionID = makeSale(t, c, 1.1, .1)

	_, err = c.bt.Transaction().Settle(transactionID)
	if err != nil {
		t.Error("failed to settle sale:", err)
	}
	_, err = c.RefundSale(transactionID)
	if err != nil {
		t.Error("failed to refund sale:", err)
	}
}

func TestMain(m *testing.M) {
	t, err := time.Parse(dateOfBirthFormat, "01-01-1989")
	if err != nil {
		log.Fatalf("failed to get dob time: %v", err)
	}
	updateSubMerchantReq := &SubMerchantInfo{
		ID:          "01234567890123456789532345678912",
		FirstName:   "Kayle",
		LastName:    "Gishen",
		Email:       "kayle@test.com",
		DateOfBirth: t,
		Address: types.Address{
			APT:     "Suite 404",
			Street:  "1 E Main St",
			City:    "Chicago",
			State:   "IL",
			Zip:     "60622",
			Country: "USA",
		},
		AccountNumber: "5836569",
		RoutingNumber: "071101307",
	}
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatalf("failed to get aetest context: %v", err)
	}
	defer done()

	c := New(ctx)
	user := new(types.User)
	subMerchantID, err = c.UpdateSubMerchant(user, updateSubMerchantReq)
	if err != nil {
		log.Fatalln("failed to create test sub merchant: %v", err)
	}
	log.Printf("using subMerchantID %s", subMerchantID)
	code := m.Run()
	os.Exit(code)
}
