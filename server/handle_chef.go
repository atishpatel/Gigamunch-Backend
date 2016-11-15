package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"

	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	bucketName          string
	errInternal         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error while uploading file."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
	projectID           string
)

type urlResp struct {
	URL string               `json:"url"`
	Err errors.ErrorWithCode `json:"err"`
}

func handleUpload(w http.ResponseWriter, req *http.Request) {
	resp := new(urlResp)
	ctx := appengine.NewContext(req)

	defer handleURLResp(ctx, w, resp)

	time.Sleep(500 * time.Millisecond) // check if failed to find blob file bug is fixed with this

	// get file
	blobs, _, err := blobstore.ParseUpload(req)
	if err != nil {
		resp.Err = errInvalidParameter.WithMessage("Error parsing multipart form.").WithError(err)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		resp.Err = errInvalidParameter.WithMessage("No file was uploaded.")
		return
	}
	opts := &image.ServingURLOptions{
		Secure: true,
		Crop:   true,
	}
	time.Sleep(500 * time.Millisecond) // check if failed to find blob file bug is fixed with this
	// ctx, _ = context.WithDeadline(ctx, time.Now().Add(60*time.Second))
	url, err := image.ServingURL(ctx, file[0].BlobKey, opts)
	if err != nil {
		deadline, _ := ctx.Deadline()
		resp.Err = errInternal.WithError(err).Wrapf("failed to get image.ServingURL (blobkey: %v) (now:%v context.Deadline:%v)", file[0].BlobKey, time.Now(), deadline)

		return
	}
	resp.URL = url.String()
}

func handleGetUploadURL(w http.ResponseWriter, req *http.Request) {
	resp := new(urlResp)
	ctx := appengine.NewContext(req)
	defer handleURLResp(ctx, w, resp)
	if bucketName == "" {
		bucketName = config.GetBucketName(ctx)
	}
	// get user
	user, err := getUserFromCookie(req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return
	}
	opts := &blobstore.UploadURLOptions{
		StorageBucket: fmt.Sprintf("%s/%s", bucketName, user.ID),
	}
	uploadURL, err := blobstore.UploadURL(ctx, "/upload", opts)
	if err != nil {
		resp.Err = errInternal.WithError(err).Wrap("error getting blobstore.UploadURL")
		return
	}
	resp.URL = uploadURL.String()
}

func handleURLResp(ctx context.Context, w http.ResponseWriter, resp *urlResp) {
	// encode json resp and log errors
	if projectID == "" {
		projectID = os.Getenv("PROJECTID")
	}
	if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
		utils.Criticalf(ctx, "SERVER: Error uploading file: %+v", resp.Err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		utils.Criticalf(ctx, "SERVER: Error encoding json: %+v", err)
	}
}
