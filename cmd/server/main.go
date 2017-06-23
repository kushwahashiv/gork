package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"

	"github.com/gork-io/gork/transformers/gateways/grpc"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// main starts program execution.
func main() {

	app := cli.NewApp()
	app.Name = "Gork Server"
	app.Usage = "Starts the server."
	app.Description = "Gork is..."
	app.Version = fmt.Sprintf("%s (commit: %s, build date: %s)", version, commit, buildDate)
	app.Action = action

	app.Flags = cliFlags
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}

// action prepares and runs default cli command.
func action(ctx *cli.Context) (err error) {

	// Parse and validate config
	config, err := NewConfigFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "config parsing failed")
	}

	// Initialize logger
	logger, err := createLogger(config)
	if err != nil {
		return errors.Wrap(err, "logger creation failed")
	}
	defer logger.Sync()

	// Initialize gateways
	grpcGateway := grpc.NewGateway()

	// Start server
	logger.Debug("Server is starting...")
	server, err := NewServer(ServerWithGateways(
		grpcGateway,
	))
	if err != nil {
		return
	}

	err = server.Start()

	// Run until a quit signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Exit gracefully
	logger.Info("Termination signal received, shutting down gracefully...")
	server.Stop()
	return
}

// createLogger creates a new logger instance for config given.
func createLogger(config *config) (logger *zap.Logger, err error) {

	var loggerConfig zap.Config
	if config.Misc.DebugMode {
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		loggerConfig = zap.NewProductionConfig()
	}
	if config.Misc.LogFormatText {
		loggerConfig.Encoding = "console"
	}

	return loggerConfig.Build()
}
