package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/database"
	"github.com/UniqueStudio/UniqueSSO/middleware"
	"github.com/UniqueStudio/UniqueSSO/model"
	"github.com/UniqueStudio/UniqueSSO/router"
	"github.com/UniqueStudio/UniqueSSO/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xylonx/zapx"
	zapxdecoder "github.com/xylonx/zapx/decoder"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "unique-sso",
	Short: "unique studio sso service",
	PreRun: func(c *cobra.Command, args []string) {
		setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var configFilePath string

func init() {
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "./conf/conf.yaml", "path to config file")
}

func Execute() error {
	return rootCmd.Execute()
}

func setup() {
	// init zapx
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	zapx.Use(logger, zapxdecoder.OpentelemetaryDecoder)

	if err := conf.InitConf(configFilePath); err != nil {
		zapx.Error("load config from file failed", zap.String("file", configFilePath), zap.Error(err))
		os.Exit(1)
	}

	if err := database.InitDB(); err != nil {
		os.Exit(1)
	}

	if err := model.InitTables(); err != nil {
		os.Exit(1)
	}

	if err := util.SetupUtils(); err != nil {
		os.Exit(1)
	}

	if err := middleware.SetupMiddleware(); err != nil {
		os.Exit(1)
	}
}

func run() {
	// setup otel tracing
	shutdown, err := util.SetupTracing()
	defer func() {
		zapx.Info("tracing reporter is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			zapx.Error("tracing have been down down")
			return
		}
		zapx.Info("tracing reporter shut down successfully")
	}()

	if err != nil {
		zapx.Error("setup otel tracing failed", zap.Error(err))
		os.Exit(1)
	}

	/*
		Init HTTP server
	*/
	if conf.SSOConf.Application.Mode == common.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	router.InitRouter(r)

	httpAddr := conf.SSOConf.Application.Host + ":" + conf.SSOConf.Application.Port
	srv := http.Server{
		Addr:         httpAddr,
		Handler:      r,
		ReadTimeout:  time.Second * time.Duration(conf.SSOConf.Application.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(conf.SSOConf.Application.WriteTimeout),
	}

	/*
		INIT gRPC server
	*/
	grpcAddr := conf.SSOConf.Application.Host + ":" + conf.SSOConf.Application.RPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := new(grpc.Server)
	if conf.SSOConf.Application.Mode == common.DebugMode {
		creds, err := credentials.NewServerTLSFromFile(
			conf.SSOConf.Application.RPCCertFile,
			conf.SSOConf.Application.RPCKeyFile,
		)
		if err != nil {
			zapx.Fatal("new server credentials failed", zap.Error(err))
		}
		s = grpc.NewServer(
			grpc.Creds(creds),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				otelgrpc.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				otelgrpc.UnaryServerInterceptor(),
			)),
		)
	} else {
		s = grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				otelgrpc.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				otelgrpc.UnaryServerInterceptor(),
			)),
		)
	}

	router.InitRPC(s)

	/*
		Start HTTP server
	*/
	go func() {
		zapx.Info("start http server", zap.String("host", httpAddr))
		if err := srv.ListenAndServe(); err != nil {
			zapx.Error("http run error", zap.Error(err))
		}
	}()

	/*
		Start gRPC server
	*/
	go func() {
		zapx.Info("start grpc server", zap.String("host", grpcAddr))
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	/*
		Graceful Shutdown
	*/
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s.GracefulStop()

	if err := srv.Shutdown(ctx); err != nil {
		zapx.Error("shutdown http server failed", zap.Error(err))
		os.Exit(1)
	}

}
