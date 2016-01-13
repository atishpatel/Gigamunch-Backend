package databases

import (
        "testing"

        "google.golang.org/appengine"
        "google.golang.org/appengine/aetest"
)

var (
  db Database
)

func setup() {
  db = CreateDatabase(Config{
    RedisSessionDBIP:
  	RedisSessionDBPassword:
  	RedisSessionDBPoolSize:
  })
}

func TestSaveUserSession(t *testing.T) {
   ctx, done err := aetest.NewContext()
   if err != nil {
           t.Fatal(err)
   }

   defer done()
}
