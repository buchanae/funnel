package server

import (
	"github.com/golang/gddo/httputil"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/webdash"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
)

type Config struct {
	ServiceName string
	HostName    string
	HTTPPort    string
	DisableHTTPCache bool
	Logger           logger.Config
  RPC rpc.Config
  Cert string
  Key string
}

// Server represents a Funnel server. The server handles
// RPC traffic via gRPC, HTTP traffic for the TES API,
// and also serves the web dashboard.
type Server struct {
  Conf Config
	TaskServiceServer      tes.TaskServiceServer
	EventServiceServer     events.EventServiceServer
	SchedulerServiceServer pbs.SchedulerServiceServer
	Log                    *logger.Logger
}

// Serve starts the server and does not block. This will open TCP ports
// for both RPC and HTTP.
func (s *Server) Serve(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	srvOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// API auth check.
				newAuthInterceptor(s.RPC.User, s.RPC.Password),
				// Debug logs which log every request/response.
				newDebugInterceptor(s.Log),
			),
		),
	}

	// SSL/TLS configuration
	if s.Cert != "" {
		creds, err := credentials.NewServerTLSFromFile(s.Cert, s.Key)
		if err != nil {
			return err
		}
		srvOpts = append(srvOpts, grpc.Creds(creds))
	}

	// Set up RPC APIs and HTTP gateways.
	grpcServer := grpc.NewServer(srvOpts...)
	grpcMux := runtime.NewServeMux()
	runtime.OtherErrorHandler = s.handleError

	// RPC network connection
	mux := http.NewServeMux()
  s.RPC.Host = "localhost"
	conn, err := rpc.NewConn(ctx, s.RPC)
	if err != nil {
		return err
	}

	// Register TES service
	if s.TaskServiceServer != nil {
		tes.RegisterTaskServiceServer(grpcServer, s.TaskServiceServer)
		err := tes.RegisterTaskServiceHandler(ctx, grpcMux, conn)
		if err != nil {
			return err
		}
	}

	// Register Events service
	if s.EventServiceServer != nil {
		events.RegisterEventServiceServer(grpcServer, s.EventServiceServer)
	}

	// Register Scheduler RPC service
	if s.SchedulerServiceServer != nil {
		pbs.RegisterSchedulerServiceServer(grpcServer, s.SchedulerServiceServer)
		err := pbs.RegisterSchedulerServiceHandler(ctx, grpcMux, conn)
		if err != nil {
			return err
		}
	}

	// Set up web dashboard routing.
	dashmux := http.NewServeMux()
	dashmux.Handle("/", webdash.RootHandler())
	dashfs := webdash.FileServer()
	mux.Handle("/favicon.ico", dashfs)
	mux.Handle("/static/", http.StripPrefix("/static/", dashfs))

	// Content negotiation for the API vs web dashboard.
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {

		switch negotiate(req) {
		case "html":
			// HTML was requested (by the browser)
			dashmux.ServeHTTP(resp, req)
		default:
			// Set "cache-control: no-store" to disable response caching.
			// Without this, some servers (e.g. GCE) will cache a response from ListTasks, GetTask, etc.
			// which results in confusion about the stale data.
			if s.DisableHTTPCache {
				resp.Header().Set("Cache-Control", "no-store")
			}
			grpcMux.ServeHTTP(resp, req)
		}
	})

	httpServer := &http.Server{
		Addr:    ":" + s.HTTPPort,
		Handler: mux,
	}

	// Open TCP connection for RPC
	lis, err := net.Listen("tcp", s.RPC.Address())
	if err != nil {
		return err
	}

	var srverr error
	go func() {
		srverr = grpcServer.Serve(lis)
		cancel()
	}()

	go func() {
		if s.Cert != "" {
			srverr = httpServer.ListenAndServeTLS(s.Cert, s.Key)
		} else {
			srverr = httpServer.ListenAndServe()
		}
		cancel()
	}()

	s.Log.Info("Server listening",
		"httpPort", s.HTTPPort,
		"rpcAddress", s.RPC.Address,
	)

	<-ctx.Done()
	grpcServer.GracefulStop()
	httpServer.Shutdown(context.TODO())

	return srverr
}

// handleError handles errors in the HTTP stack, logging errors, stack traces,
// and returning an HTTP error code.
func (s *Server) handleError(w http.ResponseWriter, req *http.Request, err string, code int) {
	s.Log.Error("HTTP handler error", "error", err, "url", req.URL)
	http.Error(w, err, code)
}

// negotiate determines the response type based on request headers and parameters.
// Returns either "html" or "json".
func negotiate(req *http.Request) string {
	// Allow overriding the type from a URL parameter.
	// /v1/tasks?json will force a JSON response.
	q := req.URL.Query()
	if _, html := q["html"]; html {
		return "html"
	}
	if _, json := q["json"]; json {
		return "json"
	}
	// Content negotiation means that both the dashboard's HTML and the API's JSON
	// may be served at the same path.
	// In Go 1.10 we'll be able to move to a core library for this,
	// https://github.com/golang/go/issues/19307
	switch httputil.NegotiateContentType(req, []string{"text/*", "text/html"}, "text/*") {
	case "text/html":
		return "html"
	default:
		return "json"
	}
}

// Return a new interceptor function that logs all requests at the Debug level
func newDebugInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	// Return a function that is the interceptor.
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		log.Debug(
			"received: "+info.FullMethod,
			"request", req,
		)
		resp, err := handler(ctx, req)
		log.Debug(
			"responding: "+info.FullMethod,
			"resp", resp,
			"err", err,
		)
		return resp, err
	}
}
