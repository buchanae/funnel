package server

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"runtime/debug"
  "github.com/ohsu-comp-bio/funnel/webdash"
)

var log = logger.New("server")

// Server represents a Funnel server. The server handles
// RPC traffic via gRPC, HTTP traffic for the TES API,
// and also serves the web dashboard.
type Server struct {
  RPCAddress string
  HTTPPort  string
  TaskServiceServer tes.TaskServiceServer
  SchedulerServiceServer pbf.SchedulerServiceServer
  Handler http.Handler
}

func DefaultServer(db Database, conf config.Config) *Server {
	log.Debug("Server Config", "config.Config", conf)

  mux := http.NewServeMux()
  mux.Handle("/", webdash.Handler())

	// Set "cache-control: no-store" to disable response caching.
	// Without this, some servers (e.g. GCE) will cache a response from ListTasks, GetTask, etc.
	// which results in confusion about the stale data.
	if conf.DisableHTTPCache {
		mux.Handle("/v1/", noCacheHandler(mux))
	} else {
		mux.Handle("/v1/", mux)
	}

  return &Server{
    RPCAddress: conf.RPCAddress(),
    HTTPPort: conf.HTTPPort,
    TaskServiceServer: db,
    SchedulerServiceServer: db,
    Handler: mux,
  }
}

// Start starts the server and does not block. This will open TCP ports
// for both RPC and HTTP.
func (s *Server) Serve(ctx context.Context) error {

	grpcServer := grpc.NewServer()

	// Set up HTTP proxy of gRPC API
	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
    grpc.WithInsecure(),
  }
	runtime.OtherErrorHandler = handleError

	httpServer := &http.Server{
		Addr: ":" + s.HTTPPort,
    Handler: s.Handler,
	}

  // Open TCP connection for RPC
	lis, err := net.Listen("tcp", s.RPCAddress)
	if err != nil {
		return err
	}

  // Register TES service
  if s.TaskServiceServer != nil {
	  tes.RegisterTaskServiceServer(grpcServer, s.TaskServiceServer)
    err := tes.RegisterTaskServiceHandlerFromEndpoint(
      ctx, grpcMux, s.RPCAddress, opts,
    )
    if err != nil {
      return err
    }
  }

  // Register Scheduler RPC service
  if s.SchedulerServiceServer != nil {
	  pbf.RegisterSchedulerServiceServer(grpcServer, s.SchedulerServiceServer)
    err := pbf.RegisterSchedulerServiceHandlerFromEndpoint(
      ctx, grpcMux, s.RPCAddress, opts,
    )
    if err != nil {
      return err
    }
  }

	log.Info("RPC server listening", "address", s.RPCAddress)

	go func() {
		err := grpcServer.Serve(lis)
		log.Error("RPC server error", err)
	}()

	log.Info("HTTP server listening",
		"httpPort", s.HTTPPort, "rpcAddress", s.RPCAddress,
	)
	// TODO how do we handle errors returned from grpcServer.Serve()
	//      httpServer.ListenAndServe()
	go func() {
		err := httpServer.ListenAndServe()
		log.Error("HTTP server error", err)
	}()

  select {
  case <-ctx.Done():
    grpcServer.GracefulStop()
    httpServer.Shutdown(context.TODO())
  }
	return nil
}

// handleError handles errors in the HTTP stack, logging errors, stack traces,
// and returning an HTTP error code.
func handleError(w http.ResponseWriter, req *http.Request, err string, code int) {
	log.Error("HTTP handler error", "error", err, "url", req.URL)
	debug.PrintStack()
	http.Error(w, err, code)
}

// Set a cache-control header that disables response caching
// and pass through to the next mux.
func noCacheHandler(next http.Handler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(resp, req)
	}
}
