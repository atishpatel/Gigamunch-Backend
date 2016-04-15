package server

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/cloud/storage"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"
)

var (
	bucketName          string
	errInternal         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error while uploading file."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
)

// func handleGigachefApp(w http.ResponseWriter, req *http.Request) {
// 	var err error
// 	var page []byte
// 	if appengine.IsDevAppServer() {
// 		page, err = ioutil.ReadFile("chef/app/index.html")
// 	} else {
// 		page = chefIndexPage
// 	}
//
// 	if err != nil {
// 		ctx := appengine.NewContext(req)
// 		utils.Errorf(ctx, "Error reading login page: %+v", err)
// 	}
// 	w.Write(page)
// }

// func middlewareLoggedIn(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		user := CurrentUser(w, req)
// 		if user == nil {
// 			http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
// 			return
// 		}
// 		h.ServeHTTP(w, req)
// 	})
// }

func handleUpload(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var returnErr errors.ErrorWithCode
	var resp struct {
		URL string               `json:"url"`
		Err errors.ErrorWithCode `json:"err"`
	}
	ctx := appengine.NewContext(req)

	defer func() {
		// encode json resp and log errors
		resp.Err = returnErr
		if returnErr.Code != 0 && returnErr.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "Error uploading file: %+v", returnErr)
		}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			utils.Errorf(ctx, "Error decoding json: %+v", err)
		}
	}()
	// get user
	user, err := getUserFromCookie(req)
	if err != nil {
		returnErr = errors.GetErrorWithCode(err)
		return
	}
	// get file
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		returnErr = errInvalidParameter.WithMessage("File is invalid.").WithError(err)
		return
	}
	defer func() {
		if file != nil {
			err = file.Close()
			if err != nil {
				utils.Errorf(ctx, "Error closing file: %+v", err)
			}
		}
	}()
	// make sure file is an image
	if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
		returnErr = errInvalidParameter.WithMessage("Invalid file format.")
		return
	}
	// decode image
	img, err := imaging.Decode(file)
	if err != nil {
		returnErr = errInternal.WithError(err)
		return
	}
	// resize image
	img = resizeAndConvert(img)
	// generate obj name of UserID/rand
	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	objName := user.ID + "/" + id
	// save to bucket
	err = uploadToBucket(ctx, bucketName, objName, img)
	if err != nil {
		returnErr = errors.GetErrorWithCode(err)
	}
	resp.URL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objName)
}

func uploadToBucket(ctx context.Context, bucketName, objName string, img image.Image) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errInternal.WithError(err)
	}
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objName)
	wr := obj.NewWriter(ctx)
	defer func() {
		err = wr.Close()
		if err != nil {
			utils.Errorf(ctx, "Error closing file writer: %+v", err)
		}
	}()
	wr.CacheControl = "max-age=1209600" // cache for 14 days
	err = jpeg.Encode(wr, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return errInternal.WithError(err)
	}
	return nil
}

func resizeAndConvert(img image.Image) image.Image {
	// rotate images that are portrait
	if img.Bounds().Dx() < img.Bounds().Dy() {
		img = imaging.Rotate90(img)
	}
	ratio := 9.0 / 16.0
	width := img.Bounds().Dx()
	// crop
	img = imaging.CropCenter(img, width, int(float64(width)*ratio))
	// resize
	img = imaging.Resize(img, 1920, 1080, imaging.Lanczos)
	return img
}

func init() {
	// var err error
	// TODO switch to template with footer and stuff in different page
	// chefIndexPage, err = ioutil.ReadFile("chef/app/index.html")
	// if err != nil {
	// 	log.Fatal("chef/app/index.html not found")
	// }
	bucketName = "gigamunch-dev-images"
}
