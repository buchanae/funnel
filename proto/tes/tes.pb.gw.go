// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: tes.proto

/*
Package tes is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package tes

import (
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

func request_TaskService_GetServiceInfo_0(ctx context.Context, marshaler runtime.Marshaler, srv TaskServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, error) {
	var protoReq ServiceInfoRequest

	msg, err := srv.GetServiceInfo(ctx, &protoReq)
	return msg, err

}

func request_TaskService_CreateTask_0(ctx context.Context, marshaler runtime.Marshaler, srv TaskServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, error) {
	var protoReq Task

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := srv.CreateTask(ctx, &protoReq)
	return msg, err

}

var (
	filter_TaskService_ListTasks_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_TaskService_ListTasks_0(ctx context.Context, marshaler runtime.Marshaler, srv TaskServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, error) {
	var protoReq ListTasksRequest

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_TaskService_ListTasks_0); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := srv.ListTasks(ctx, &protoReq)
	return msg, err

}

var (
	filter_TaskService_GetTask_0 = &utilities.DoubleArray{Encoding: map[string]int{"id": 0}, Base: []int{1, 1, 0}, Check: []int{0, 1, 2}}
)

func request_TaskService_GetTask_0(ctx context.Context, marshaler runtime.Marshaler, srv TaskServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, error) {
	var protoReq GetTaskRequest

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_TaskService_GetTask_0); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := srv.GetTask(ctx, &protoReq)
	return msg, err

}

func request_TaskService_CancelTask_0(ctx context.Context, marshaler runtime.Marshaler, srv TaskServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, error) {
	var protoReq CancelTaskRequest

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	msg, err := srv.CancelTask(ctx, &protoReq)
	return msg, err

}

// RegisterTaskServiceHandler registers the http handlers for service TaskService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterTaskServiceHandler(mux *runtime.ServeMux) error {

	mux.Handle("GET", pattern_TaskService_GetServiceInfo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
    ctx := req.Context()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_TaskService_GetServiceInfo_0(rctx, inboundMarshaler, srv, req, pathParams)

		forward_TaskService_GetServiceInfo_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_TaskService_CreateTask_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
    ctx := req.Context()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_TaskService_CreateTask_0(rctx, inboundMarshaler, srv, req, pathParams)

		forward_TaskService_CreateTask_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_TaskService_ListTasks_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_TaskService_ListTasks_0(rctx, inboundMarshaler, srv, req, pathParams)

		forward_TaskService_ListTasks_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_TaskService_GetTask_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_TaskService_GetTask_0(rctx, inboundMarshaler, srv, req, pathParams)

		forward_TaskService_GetTask_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_TaskService_CancelTask_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_TaskService_CancelTask_0(rctx, inboundMarshaler, srv, req, pathParams)

		forward_TaskService_CancelTask_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_TaskService_GetServiceInfo_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "tasks", "service-info"}, ""))

	pattern_TaskService_CreateTask_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "tasks"}, ""))

	pattern_TaskService_ListTasks_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "tasks"}, ""))

	pattern_TaskService_GetTask_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"v1", "tasks", "id"}, ""))

	pattern_TaskService_CancelTask_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"v1", "tasks", "id"}, "cancel"))
)

var (
	forward_TaskService_GetServiceInfo_0 = runtime.ForwardResponseMessage

	forward_TaskService_CreateTask_0 = runtime.ForwardResponseMessage

	forward_TaskService_ListTasks_0 = runtime.ForwardResponseMessage

	forward_TaskService_GetTask_0 = runtime.ForwardResponseMessage

	forward_TaskService_CancelTask_0 = runtime.ForwardResponseMessage
)
