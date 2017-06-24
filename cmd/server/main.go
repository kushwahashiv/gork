package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"

	"net"

	"github.com/go-redis/redis"
	"github.com/gork-io/gork/services/resources"
	"github.com/gork-io/gork/transformers/gateways/grpc"
	"github.com/gork-io/gork/transformers/gateways/grpc/controllers"
	redis_repo "github.com/gork-io/gork/transformers/repositories/redis"
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

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Db.Redis.Hostname, config.Db.Redis.Port),
		DB:       config.Db.Redis.Database,
		Password: config.Db.Redis.Password,
	})

	// Initialize repositories
	queuesRepo := redis_repo.NewQueuesRepository(redisClient)

	// Initialize services
	queuesSvc := resources.NewQueues(queuesRepo)

	// Initialize gateways
	listener, err := net.Listen("tcp", net.JoinHostPort(config.Gtw.Grpc.Hostname, config.Gtw.Grpc.Port))
	if err != nil {
		return
	}
	grpcGateway := grpc.NewGateway(listener, controllers.NewQueues(queuesSvc))

	// Initialize server
	server, err := NewServer(ServerWithGateways(grpcGateway))
	if err != nil {
		return
	}

	// Start server
	logger.Debug("Server is starting...")
	err = server.Start()
	if err != nil {
		return
	}

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
