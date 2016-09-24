package like

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
)

func TestLike(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	c := New(ctx)
	userID := randUserID()
	itemID := randItemID()
	// test
	err = c.Like(userID, itemID)
	if err != nil {
		t.Fatal("Error while liking item: ", err)
	}
	rows, err := mysqlDB.Query("SELECT item_id FROM `like` WHERE user_id=? AND item_id=?", userID, itemID)
	if err != nil {
		t.Fatal("Error selecting liked item: ", err)
	}
	defer handleTestCloser(t, rows)
	var tmpItemID int64
	for rows.Next() {
		_ = rows.Scan(&tmpItemID)
	}
	if tmpItemID != itemID {
		t.Fatal("ItemID isn't the same: ", tmpItemID, itemID)
	}
}

func randUserID() string {
	return strconv.FormatInt(int64(time.Now().Nanosecond()), 36)
}

func randItemID() int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Int63()
}

func TestUnlike(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	c := New(ctx)
	userID := randUserID()
	itemID := randItemID()
	err = c.Like(userID, itemID)
	if err != nil {
		t.Fatal("Error while liking item: ", err)
	}
	// test
	err = c.Unlike(userID, itemID)
	if err != nil {
		t.Fatal("Error while unliking item: ", err)
	}
	rows, err := mysqlDB.Query("SELECT item_id FROM `like` WHERE user_id=? AND item_id=?", userID, itemID)
	if err != nil {
		t.Fatal("Error selecting liked item: ", err)
	}
	defer handleTestCloser(t, rows)
	if rows.Next() {
		t.Fatal("unlike item not deleted. ")
	}
}

func TestGetNumLikes(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	c := New(ctx)
	userID := randUserID()
	repeatedItemID := randItemID()
	items := []int64{repeatedItemID, randItemID(), randItemID(), repeatedItemID}
	for _, item := range items {
		err = c.Like(userID, item)
		if err != nil {
			t.Fatal("Error while liking item: ", err)
		}
	}
	err = c.Like(randUserID(), items[0])
	if err != nil {
		t.Fatal("Error while liking item: ", err)
	}
	err = c.Unlike(userID, items[1])
	if err != nil {
		t.Fatal("Error while unliking item: ", err)
	}
	// test
	wantNumLikes := []int{2, 0, 1, 2}
	numLikes, err := c.GetNumLikes(items)
	if err != nil {
		t.Fatal("Error while calling LikesItems: ", err)
	}
	if len(wantNumLikes) != len(numLikes) {
		t.Fatal("Array size of want and actual is different")
	}
	for i := range numLikes {
		if numLikes[i] != wantNumLikes[i] {
			t.Fatalf("wantResults(%v) and actualResults(%v) aren't the same.", wantNumLikes, numLikes)
		}
	}
}

func TestLikesItems(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	c := New(ctx)
	userID := randUserID()
	repeatedItemID := randItemID()
	items := []int64{repeatedItemID, randItemID(), randItemID(), repeatedItemID}
	for _, item := range items {
		err = c.Like(userID, item)
		if err != nil {
			t.Fatal("Error while liking item: ", err)
		}
	}
	err = c.Like(randUserID(), items[0])
	if err != nil {
		t.Fatal("Error while liking item: ", err)
	}
	err = c.Unlike(userID, items[1])
	if err != nil {
		t.Fatal("Error while unliking item: ", err)
	}
	// test
	wantLikes := []bool{true, false, true, true}
	wantNumLikes := []int{2, 0, 1, 2}
	likes, numLikes, err := c.LikesItems(userID, items)
	if err != nil {
		t.Fatal("Error while calling LikesItems: ", err)
	}
	if len(wantLikes) != len(likes) || len(wantNumLikes) != len(numLikes) {
		t.Fatal("Array size of want and actual is different")
	}
	for i := range likes {
		if likes[i] != wantLikes[i] {
			t.Fatalf("wantResults(%v) and actualResults(%v) aren't the same.", wantLikes, likes)
		}
	}
	for i := range numLikes {
		if numLikes[i] != wantNumLikes[i] {
			t.Fatalf("wantResults(%v) and actualResults(%v) aren't the same.", wantNumLikes, numLikes)
		}
	}
}

func handleTestCloser(t *testing.T, c closer) {
	err := c.Close()
	if err != nil {
		t.Error("Error while closing rows: ", err)
	}
}
