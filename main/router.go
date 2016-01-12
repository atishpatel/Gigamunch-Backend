package main

import (
	// golang
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// appengine
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	// dev
	aelog "google.golang.org/appengine/log"
)

const (
	loginURL     = "/login"
	gitkitURL    = "/gitkit"
	oobActionURL = "/sendEmail"
	signOutURL   = "/signOut"

	//TODO(Atish): Remove after ComingSoon Phase.
	comingSoonURL = "/"
	thankyouURL   = "/thankyou"
)

func init() {
	// TODO(Atish): remove after ComingSoon Phase
	http.HandleFunc(thankyouURL, handleThankYou)
	var err error
	thankyouPage, err = ioutil.ReadFile("app/thankyou.html")
	if err != nil {
		log.Fatalf("Failed to read thankyou page %#v", err)
	}
}

/*
 * ComingSoon stuff
 * TODO(Atish): remove after ComingSoon Phase
 */

var thankyouPage []byte

type ComingSoon struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date" datastore:",noindex"`
	Name  string    `json:"name"`
}

func handleThankYou(w http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	emailAddress := req.FormValue("email")
	name := req.FormValue("name")
	terp := req.FormValue("terp")
	aelog.Infof(c, "email: %s, name: %s, terp: %s ", emailAddress, name, terp)
	if terp != "" {
		return
	}
	if emailAddress == "" {
		aelog.Infof(c, "No email address. ")
		w.Write(thankyouPage)
		return
	}
	key := datastore.NewKey(c, "ComingSoon", emailAddress, 0, nil)
	entry := &ComingSoon{}
	err := datastore.Get(c, key, entry)
	if err == datastore.ErrNoSuchEntity {
		entry.Date = time.Now()
		entry.Email = emailAddress
		entry.Name = name
		_, err = datastore.Put(c, key, entry)
		if err != nil {
			aelog.Errorf(c, "Error putting ComingSoonEmail in datastore ", err)
		}
		w.Write(thankyouPage)
		return
	}
	aelog.Errorf(c, "Error email already registered ComingSoonEmail: emailaddress - %s, err - %#v", emailAddress, err)
	w.Write(thankyouPage)
}
