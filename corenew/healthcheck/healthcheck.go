package healthcheck

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

const (
	kindDevice = "Device"
)

var (
	managerNumbers = [...]string{"9316445311", "6155454989", "6153975516", "9316446755"}
	errDatastore   = errors.InternalServerError
	errBadRequest  = errors.BadRequestError
)

// Device is a device that is reporting in for healthcheck.
type Device struct {
	ID            string    `json:"id"`
	Disable       bool      `json:"disable"`
	NumAlerts     int       `json:"num_alerts"`
	IsPowerSensor bool      `json:"is_power_sensor"`
	LastCheckin   time.Time `json:"last_checkin"`
}

// Client is the client for this package.
type Client struct {
	ctx context.Context
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	return &Client{
		ctx: ctx,
	}
}

// Checkin updates the device.
func (c *Client) Checkin(d *Device) error {
	if d.ID == "" {
		return errBadRequest.Annotate("id cannot be empty")
	}
	oldDevice, err := get(c.ctx, d.ID)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to get")
	}
	if oldDevice.Disable {
		msg := fmt.Sprintf("Baby I'm back! The following device is back online:\n %+v", d.ID)
		messageC := message.New(c.ctx)
		numbers := managerNumbers
		for _, number := range numbers {
			err = messageC.SendDeliverySMS(number, msg)
			if err != nil {
				logging.Errorf(c.ctx, "failed to message.SendAdminSMS: %+v", err)
			}
		}
	}

	err = put(c.ctx, d.ID, d)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	return nil
}

// CheckPowerSensors checks all the PowerSensors.
func (c *Client) CheckPowerSensors() error {
	devices, err := getPowerSensors(c.ctx)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to getPowerSensors")
	}
	if len(devices) == 0 {
		return nil
	}
	var failingDevices []string
	for _, d := range devices {
		if !d.Disable && time.Since(d.LastCheckin) > time.Minute*1 {
			failingDevices = append(failingDevices, d.ID)
			d.NumAlerts++
			if d.NumAlerts >= 3 {
				d.Disable = true
				d.NumAlerts = 0
			}
			err = put(c.ctx, d.ID, d)
			if err != nil {
				logging.Errorf(c.ctx, "failed to put: %+v", err)
			}
		}
	}
	if len(failingDevices) != 0 {
		msg := fmt.Sprintf("ALERT! POWER OUTAGE! The following devices don't have power:\n %+v", failingDevices)
		messageC := message.New(c.ctx)
		numbers := managerNumbers
		for _, number := range numbers {
			err = messageC.SendDeliverySMS(number, msg)
			if err != nil {
				logging.Errorf(c.ctx, "failed to message.SendAdminSMS: %+v", err)
			}
		}
	}
	return nil
}

func get(ctx context.Context, id string) (*Device, error) {
	i := new(Device)
	key := datastore.NewKey(ctx, kindDevice, id, 0, nil)
	err := datastore.Get(ctx, key, i)
	return i, err
}

func getPowerSensors(ctx context.Context) ([]*Device, error) {
	var dst []*Device
	var err error
	_, err = datastore.NewQuery(kindDevice).Filter("IsPowerSensor=", true).GetAll(ctx, &dst)
	if err != nil && err.Error() != "(0 errors)" {
		return nil, err
	}
	return dst, nil
}

func put(ctx context.Context, id string, i *Device) error {
	var err error
	key := datastore.NewKey(ctx, kindDevice, id, 0, nil)
	_, err = datastore.Put(ctx, key, i)
	return err
}
