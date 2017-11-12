package events

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"time"
)

// RPCWriter is a type which writes Events to RPC.
type RPCWriter struct {
	client        EventServiceClient
	updateTimeout time.Duration
}

// NewRPCWriter returns a new RPCWriter instance.
func NewRPCWriter(ctx context.Context, conf rpc.Config) (*RPCWriter, error) {
	conn, err := rpc.NewConn(ctx, conf)
	if err != nil {
		return nil, err
	}
	cli := NewEventServiceClient(conn)

	return &RPCWriter{cli, conf.Timeout}, nil
}

// Write writes the event. The RPC call may timeout, based on the timeout given
// by the configuration in NewRPCWriter.
func (r *RPCWriter) Write(e *Event) error {
	ctx, cleanup := context.WithTimeout(context.Background(), r.updateTimeout)
	_, err := r.client.CreateEvent(ctx, e)
	cleanup()
	return err
}

// Close closes the writer.
func (r *RPCWriter) Close() error {
	return nil
}
