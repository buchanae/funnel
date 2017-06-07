package logutils

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"google.golang.org/grpc/grpclog"
	"os"
)

// Wrap our logger to fit the grpc logger interface
type grpclogger struct {
}

func (g *grpclogger) Fatal(args ...interface{}) {
	logger.Error("grpc", "msg", fmt.Sprint(args))
}
func (g *grpclogger) Fatalf(format string, args ...interface{}) {
	logger.Error("grpc", "msg", fmt.Sprint(args))
}
func (g *grpclogger) Fatalln(args ...interface{}) {
	logger.Error("grpc", "msg", fmt.Sprint(args))
}
func (g *grpclogger) Print(args ...interface{}) {
	logger.Info("grpc", "msg", fmt.Sprint(args))
}
func (g *grpclogger) Printf(format string, args ...interface{}) {
	logger.Info("grpc", "msg", fmt.Sprint(args))
}
func (g *grpclogger) Println(args ...interface{}) {
	logger.Info("grpc", "msg", fmt.Sprint(args))
}

func init() {
	// grpclog says to only call this from init(), so here we are
	grpclog.SetLogger(&grpclogger{})
}

// Configure configures the logging level and output path.
func Configure(conf config.Config) {
	logger.SetLevel(conf.LogLevel)
	logger.DisableTimestamp(!conf.TimestampLogs)

	// TODO Good defaults, configuration, and reusable way to configure logging.
	//      Also, how do we get this to default to /var/log/tes/worker.log
	//      without having file permission problems? syslog?
	if conf.LogPath != "" {
		logFile, err := os.OpenFile(
			conf.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
		)
		if err != nil {
			logger.Error("Can't open log output file", "path", conf.LogPath)
		} else {
			logger.SetOutput(logFile)
		}
	}
}
