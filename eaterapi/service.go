package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"google.golang.org/appengine"
	"google.golang.org/grpc"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/eater"
)

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