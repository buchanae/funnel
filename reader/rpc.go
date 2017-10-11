package reader

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"google.golang.org/grpc"
	"time"
)

// RPCReader provides read access to tasks from the funnel server over gRPC.
type RPCReader struct {
	cli tes.TaskServiceClient
}

// NewRPCReader returns a new TES RPC client with the given configuration,
// including transport security, basic password auth, and dial timeout.
func NewRPCReader(address, password string) (*RPCReader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
    address,
		grpc.WithInsecure(),
		util.PerRPCPassword(password),
	)
	if err != nil {
		return nil, err
	}
	cli := tes.NewTaskServiceClient(conn)

	return &RPCReader{cli}, nil
}
