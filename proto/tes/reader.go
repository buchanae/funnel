package tes

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ReadOnlyServer interface {
	ListTasks(ctx context.Context, in *ListTasksRequest) (*ListTasksResponse, error)
	GetTask(ctx context.Context, in *GetTaskRequest) (*Task, error)
}

type ReadOnlyClient interface {
	ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponse, error)
	GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*Task, error)
}
