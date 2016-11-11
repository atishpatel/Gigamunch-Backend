package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"

	// mysql "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
	"google.golang.org/grpc"

	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

func getGRPCError(err error, detail string) *pb.Error {
	cErr := errors.Wrap(detail, err)
	return &pb.Error{
		Code:    cErr.Code,
		Message: cErr.Message,
		Detail:  cErr.Detail,
	}
}

func handleResp(ctx context.Context, fnName string, err *pb.Error) {
	if err == nil { // there was no error
		return
	}
	code := err.Code
	if code == errors.CodeInvalidParameter {
		utils.Warningf(ctx, "%s invalid parameter: %+v", fnName, *err)
		return
	} else if code != 0 {
		utils.Errorf(ctx, "%s err: %+v", fnName, *err)
	}
}

// processErrorChans returns an error if any of the error channels return an error
func processErrorChans(errs ...<-chan error) error {
	var err error
	for _, v := range errs {
		err = <-v
		if err != nil {
			return err
		}
	}
	return nil
}

type service struct{}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEaterServer(s, &service{})
	go func() {
		http.HandleFunc("/_ah/health", healthCheckHandler)
		http.HandleFunc("/", fronthandler)
		http.HandleFunc(tasks.ProcessInquiryURL, handleProcessInquiry)
		// http.HandleFunc("/sql", sqlTest)
		appengine.Main()
	}()
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to s.Serve: %+v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func fronthandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Whatcha doin' here?")
}

// func sqlTest(w http.ResponseWriter, r *http.Request) {
// 	v := r.URL.Query()
// 	q := v.Get("cs")
// 	output := fmt.Sprintf("running sqlConnection test %s \n", q)
// 	mysqlDB, err := sql.Open("mysql", q)
// 	if err != nil {
// 		output += "Error doing sql.Open: " + err.Error()
// 	}
// 	err = mysqlDB.Ping()
// 	if err != nil {
// 		output += "Error doing sql.Ping: " + err.Error()
// 	}
// 	fmt.Fprintln(w, output)
// }

func handleProcessInquiry(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	inquiryID, err := tasks.ParseInquiryID(req)
	if err != nil {
		utils.Errorf(ctx, "Failed to parse process inquiry request. Err: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inquiryC := inquiry.New(ctx)
	err = inquiryC.Process(inquiryID)
	if err != nil {
		utils.Criticalf(ctx, "Failed to process inquiry(%d). Err: %v", inquiryID, err)
		taskC := tasks.New(ctx)
		err = taskC.AddProcessInquiry(inquiryID, time.Now().Add(1*time.Hour))
		if err != nil {
			utils.Criticalf(ctx, "Failed to add inquiry(%d) in processInquiry queue. Err: %v", inquiryID, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
