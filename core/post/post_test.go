package post

import (
	"testing"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/misc/testhelper"

	"google.golang.org/appengine/aetest"
)

func TestGetClosedPosts(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	connectSQL(ctx)
	postIDs := []int64{testhelper.RandInt(), testhelper.RandInt(), testhelper.RandInt()}
	p := new(Post)
	closeTime := time.Now()
	for i, id := range postIDs {
		p.ClosingDateTime = closeTime.Add(time.Duration(i) * time.Minute)
		err = insertLivePost(ctx, id, p)
		if err != nil {
			t.Fatal("error inserting live post: ", err)
		}
	}

	// test
	postC := New(ctx)
	ids, _, err := postC.GetClosedPosts()
	if err != nil {
		t.Fatal("error while getting closed posts: ", err)
	}
	if len(ids) != 1 {
		t.Fatal("returned more than 1 closed ids")
	}
	if ids[0] != postIDs[0] {
		t.Fatalf("want postID(%d). received id(%d)", postIDs[0], ids[0])
	}
	// clean up
	for _, id := range postIDs {
		err = removeLivePost(ctx, id)
		if err != nil {
			t.Fatal("error deleteing posts ", err)
		}
	}
}
