package util

import (
	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/pb/sms"

	"github.com/xylonx/zapx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var OpenClient sms.SMSServiceClient

func SetupUtils() error {
	if err := setupOpenPlatformGrpc(); err != nil {
		zapx.Error("set up open platform grpc client failed", zap.Error(err))
		return err
	}
	return nil
}

func setupOpenPlatformGrpc() error {
	c, err := grpc.Dial(
		conf.SSOConf.OpenPlatform.GrpcAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		zapx.Error("dial to open platform failed", zap.Error(err))
		return err
	}

	OpenClient = sms.NewSMSServiceClient(c)
	return nil
}
