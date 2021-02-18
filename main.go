package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	sdk "proto_demo/sdk" // Update
	"strings"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcEndpoint    = flag.String("grpc-endpoint", ":9090", "gRPC server endpoint")
	grpcWebEndpoint = flag.String("grpc-web-endpoint", ":9091", "gRPC web endpoint")
	grpcGWEndpoint  = flag.String("grpc-gateway-endpoint", ":9092", "gRPC gateway endpoint")
)

type userServer struct {
	sdk.UnimplementedUserServiceServer
}

func (u userServer) GetUser(ctx context.Context, empty *emptypb.Empty) (*sdk.User, error) {
	return &sdk.User{FullName: "hello world"}, nil
}

func (u userServer) AddUser(ctx context.Context, user *sdk.User) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

type OrgServer struct {
	sdk.UnimplementedOrganizationServiceServer
}

func (u OrgServer) GetOrganizations(ctx context.Context, empty *emptypb.Empty) (*sdk.GetOrganizationsResp, error) {

	return &sdk.GetOrganizationsResp{List: []*sdk.Organization{{Alias: "best organ"}}}, nil
}

func (u OrgServer) AddOrganizationByProvider(ctx context.Context, req *sdk.AddOrganizationByProviderReq) (*sdk.Organization, error) {
	return &sdk.Organization{Alias: "organ-1"}, nil
}

func run() error {
	e := echo.New()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userSrv := &userServer{}
	orgSrv := &OrgServer{}
	lis, err := net.Listen("tcp", *grpcEndpoint)
	noError(err)
	//opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcSrv := grpc.NewServer()
	sdk.RegisterUserServiceServer(grpcSrv, userSrv)
	sdk.RegisterOrganizationServiceServer(grpcSrv, orgSrv)

	go func() {
		err := grpcSrv.Serve(lis)
		noError(err)
	}()
	wrappedServer := grpcweb.WrapServer(grpcSrv,
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}))

	mux1 := runtime.NewServeMux()
	err = sdk.RegisterUserServiceHandlerServer(ctx, mux1, userSrv)
	noError(err)
	e.Any("/some/prefix/*", func(c echo.Context) error {
		req := c.Request()
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/some/prefix")
		mux1.ServeHTTP(c.Response(), req)
		wrappedServer.ServeHTTP(c.Response(), req)
		return nil
	})
	//err := sdk.RegisterUserServiceHandlerFromEndpoint(ctx, mux1, *grpcServerEndpoint, opts)
	//if err != nil {
	//	return err
	//}

	mux2 := runtime.NewServeMux()
	err = sdk.RegisterOrganizationServiceHandlerServer(ctx, mux2, orgSrv)
	noError(err)
	e.Any("/another/prefix/*", func(c echo.Context) error {
		req := c.Request()
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/another/prefix")
		mux2.ServeHTTP(c.Response(), req)
		wrappedServer.ServeHTTP(c.Response(), req)
		return nil
	})

	// Start HTTP server (and proxy calls to gRPC server endpoint)

	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedServer.IsGrpcWebRequest(req) {
			wrappedServer.ServeHTTP(resp, req)
		} else {
			// Fall back to other servers.
			e.ServeHTTP(resp, req)
		}
		//http.DefaultServeMux.ServeHTTP(resp, req)
	})
	webSrv := &http.Server{
		Addr:    *grpcWebEndpoint,
		Handler: handler,
	}

	return webSrv.ListenAndServe()
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func noError(err error) {
	if err != nil {
		glog.Fatal(err)
	}
}
