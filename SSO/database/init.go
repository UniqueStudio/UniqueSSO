package database

import (
	"context"
	"time"

	"github.com/UniqueStudio/UniqueSSO/conf"

	"github.com/go-redis/redis/v8"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func InitDB() (err error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: conf.SSOConf.Database.PostgresDSN,
	}))
	if err != nil {
		zapx.Error("open database failed", zap.Error(err))
		return err
	}
	DB = db

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	RedisClient, err = initRedis(ctx)
	if err != nil {
		return err
	}

	return nil
}

func initRedis(ctx context.Context) (*redis.Client, error) {
	// init redis
	rclient := redis.NewClient(&redis.Options{
		Addr:     conf.SSOConf.Redis.Addr,
		Password: conf.SSOConf.Redis.Password,
		DB:       conf.SSOConf.Redis.DB,
	})
	pong, err := rclient.Ping(ctx).Result()
	if err != nil {
		zapx.Error("ping redis error", zap.Error(err))
		return nil, err
	}
	zapx.Info("ping redis success", zap.String("result", pong))
	return rclient, nil
}
