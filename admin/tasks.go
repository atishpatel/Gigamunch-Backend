package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func setupTasksHandlers() {
	http.HandleFunc("/admin/task/SetupTags", handler(SetupTags))
	http.HandleFunc("/admin/task/CheckPowerSensors", handler(CheckPowerSensors))
}

// SetupTags sets up tags for culture preview email and culture email 2 weeks in advance.
func SetupTags(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	nextCultureDate := time.Now().Add(time.Hour * 7 * 24)
	for nextCultureDate.Weekday() != time.Monday {
		nextCultureDate = nextCultureDate.Add(24 * time.Hour)
	}
	nextPreviewDate := nextCultureDate
	mailC, err := mail.NewClient(ctx, log)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.NewClient")
	}
	mailReq := &mail.UserFields{
		Email:   "atish@gigamunchapp.com",
		AddTags: []mail.Tag{mail.GetCultureEmailTag(nextCultureDate), mail.GetPreviewEmailTag(nextPreviewDate)},
	}
	err = mailC.UpdateUser(mailReq)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to mail.UpdateUser")
	}
	return errors.NoError
}

// CheckPowerSensors checks all the PowerSensors.
func CheckPowerSensors(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	healthC := healthcheck.New(ctx)
	err = healthC.CheckPowerSensors()
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to health.CheckPowerSensors")
	}
	return errors.NoError
}
