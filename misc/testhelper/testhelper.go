package testhelper

import (
	"github.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

const (
	FoodPhotoURL   = "http://lorempixel.com/1920/1080/food"
	PersonPhotoURL = "http://lorempixel.com/1920/1080/people"
)

func GetGigamuncher(ctx context.Context) (string, gigamuncher.Gigamuncher) {
	// TODO create datastore if not exist
	return "gigamuncher", gigamuncher.Gigamuncher{
		UserDetail: types.UserDetail{
			Name:     "Muncher Name",
			Email:    "muncher@test.com",
			PhotoURL: FoodPhotoURL,
		},
	}
}

func GetGigamuncherAddress() types.Address {
	return getAddress()
}

func GetGigachefUser() types.User {
	return types.User{
		ID:          "gigachef",
		Name:        "Chef Name",
		Email:       "chef@test.com",
		ProviderID:  "google.com",
		PhotoURL:    PersonPhotoURL,
		Permissions: 2147483647,
	}
}

func getAddress() types.Address {
	// Nashville Parthenon
	return types.Address{
		Street:  "2500 West End Avenue",
		City:    "Nashville",
		Zip:     "37203",
		State:   "TN",
		Country: "USA",
		GeoPoint: types.GeoPoint{
			Latitude:  36.149535,
			Longitude: -86.813154,
		},
	}
}
