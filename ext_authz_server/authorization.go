package main

import (
	"context"
	"fmt"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ExtAuthzServer struct {
	*authv3.UnimplementedAuthorizationServer
}

func (s *ExtAuthzServer) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	fmt.Println("Authorization Check Starts...")
	// the http headers from actual request dispatched by istio is stored in http.Headers, although it's a rpc call itself
	headers := request.GetAttributes().GetRequest().GetHttp().GetHeaders()
	if len(headers) == 0 {
		fmt.Println("Authorization permission is denied")
		return s.deny(typev3.StatusCode_Forbidden), nil
	}
	fmt.Println("Authorization permission is allowed")
	return s.allow(typev3.StatusCode_OK), nil
}

func (s *ExtAuthzServer) allow(code typev3.StatusCode) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{
			// only when Status.Code == 0 means check pass, request is allowed.
			Code: EnvoyHttpStatusCodeToGrpcCode(code),
		},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{},
		},
	}
}

func (s *ExtAuthzServer) deny(code typev3.StatusCode) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{
			Code: EnvoyHttpStatusCodeToGrpcCode(code),
		},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: code},
				Body:   fmt.Sprintf(`{"code": %s, "message": "request denied"}`, code),
			},
		},
	}
}

func EnvoyHttpStatusCodeToGrpcCode(code typev3.StatusCode) int32 {
	// grpc code definition: https://pkg.go.dev/google.golang.org/grpc/codes
	// not filled completed, it depends on which http code you want to map to grpc code
	switch code {
	case typev3.StatusCode_OK, typev3.StatusCode_Created, typev3.StatusCode_Accepted,
		typev3.StatusCode_NonAuthoritativeInformation, typev3.StatusCode_NoContent,
		typev3.StatusCode_ResetContent, typev3.StatusCode_PartialContent,
		typev3.StatusCode_MultiStatus, typev3.StatusCode_AlreadyReported:
		return int32(codes.OK)
	case typev3.StatusCode_BadRequest:
		return int32(codes.InvalidArgument)
	case typev3.StatusCode_Unauthorized:
		return int32(codes.Unauthenticated)
	case typev3.StatusCode_Forbidden:
		return int32(codes.PermissionDenied)
	case typev3.StatusCode_NotFound:
		return int32(codes.NotFound)
	case typev3.StatusCode_RequestTimeout:
		return int32(codes.DeadlineExceeded)
	default:
		return int32(codes.Unknown)
	}
}

func main() {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	authv3.RegisterAuthorizationServer(s, &ExtAuthzServer{})

	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
