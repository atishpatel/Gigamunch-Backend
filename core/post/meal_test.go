package post

import (
	"reflect"
	"testing"

	"github.com/atishpatel/Gigamunch-Backend/types"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestGetMultiMeal(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// setup

	testCases := []struct {
		description string
		mealIDs     []int64
		liveMeals   []*Post
		output      error
	}{
		{
			description: "Empty meal id input",
			mealIDs:     make([]int64, 0),
			liveMeals:   make([]*Post, 0),
			output:      nil,
		},
		{
			description: "Entity does not exist in Datastore",
			mealIDs:     []int64{1, 2},
			liveMeals:   make([]*Post, 2),
			output:      appengine.MultiError{datastore.ErrNoSuchEntity, nil},
		},
		{
			description: "Gets two valid Entities form Datastore",
			mealIDs:     []int64{2, 3},
			liveMeals:   make([]*Post, 2),
			output:      nil,
		},
	}
	// put stuff in Datastore
	meal2 := &types.Meal{Title: "test meal2"}
	meal3 := &types.Meal{Title: "test meal3"}
	expectedSuccessMeal := []*types.Meal{meal2, meal3}
	meal2Key := datastore.NewKey(ctx, types.KindMeal, "", 2, nil)
	meal3Key := datastore.NewKey(ctx, types.KindMeal, "", 3, nil)
	_, err = datastore.PutMulti(ctx, []*datastore.Key{meal2Key, meal3Key}, expectedSuccessMeal)
	if err != nil {
		t.Fatal("Failed to setup with put meal ", err)
	}
	// run test
	for _, test := range testCases {
		err = getMultiMeal(ctx, test.mealIDs, test.liveMeals)
		if !reflect.DeepEqual(err, test.output) {
			t.Errorf("Failed test %s | expected error: %+v | got error: %+v", test.description, test.output, err)
		}
		if err == nil {
			for i := range test.liveMeals {
				if reflect.DeepEqual(test.liveMeals[i], expectedSuccessMeal[i]) {
					t.Errorf("Failed test %s | expected meal: %+v | got meal: %+v", test.description, expectedSuccessMeal[i], test.liveMeals[i])
				}
			}
		}
	}
}
