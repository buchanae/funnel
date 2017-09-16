package rpc

import (
  "context"
	"google.golang.org/grpc"
  "time"
)

func Dial(conf Config) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
  return DialContext(ctx, conf)
}

func DialContext(ctx context.Context, conf Config) (*grpc.ClientConn, error) {
  return grpc.DialContext(ctx,
    conf.Address(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
    PerRPCPassword(conf.Password),
	)
}
