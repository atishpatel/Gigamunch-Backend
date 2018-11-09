package sub

import (
	"strings"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
)

func (c *Client) getByIDOrEmail(idOrEmail string) (*subold.Subscriber, error) {
	if strings.Contains(idOrEmail, "@") {
		return c.getByEmail(idOrEmail)
	}
	key := c.db.NameKey(c.ctx, kind, idOrEmail)
	sub := new(subold.Subscriber)
	err := c.db.Get(c.ctx, key, sub)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return sub, nil
}

func (c *Client) getByEmail(email string) (*subold.Subscriber, error) {
	var results []*subold.Subscriber
	keys, err := c.db.QueryFilter(c.ctx, kind, 0, 100, "EmailPrefs.Email=", email, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, c.db.ErrNoSuchEntity()
	}
	results[0].ID = keys[0].NameID()
	return results[0], nil
}

// returns sub, nil if found
// returns nil, nil if not found
func (c *Client) getByAddress(address *common.Address) (*subold.Subscriber, error) {
	var results []*subold.Subscriber
	_, err := c.db.QueryFilter(c.ctx, kind, 0, 300, "Address.Zip=", address.Zip, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	for _, r := range results {
		distanceInMiles := address.GreatCircleDistance(r.Address.GeoPoint)
		if distanceInMiles < .003 { // < 15 feet
			// suspect but could be apartment complex
			apt1 := strings.TrimSpace(strings.ToLower(r.Address.APT))
			apt2 := strings.TrimSpace(strings.ToLower(address.APT))
			if apt1 == apt2 {
				// matches
				return r, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) put(id string, sub *subold.Subscriber) error {
	sub.Address.Street = strings.Title(sub.Address.Street)
	sub.Address.City = strings.Title(sub.Address.City)
	sub.Address.Zip = strings.TrimSpace(sub.Address.Zip)
	var key common.Key
	var err error
	if id == "" {
		key = c.db.IncompleteKey(c.ctx, kind)
	} else {
		key = c.db.NameKey(c.ctx, kind, id)
	}
	key, err = c.db.Put(c.ctx, key, sub)
	if err != nil {
		return err
	}
	if id == "" {
		sub.ID = key.NameID()
		_, err = c.db.Put(c.ctx, key, sub)
		if err != nil {
			return err
		}
	}
	return nil
}
