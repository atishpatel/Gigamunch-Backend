package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

func setupTasksHandlers() {
	http.HandleFunc("/admin/task/SetupTags", handler(SetupTags))
}

// SetupTags sets up tags for culture preview email and culture email 2 weeks in advance.
func SetupTags(ctx context.Context, r *http.Request, log *logging.Client) Response {
	var err error
	nextCultureDate := time.Now().Add(time.Hour * 7 * 24)
	for nextCultureDate.Weekday() != time.Monday {
		nextCultureDate = nextCultureDate.Add(24 * time.Hour)
	}
	nextPreviewDate := nextCultureDate.Add(24 * time.Hour)
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
