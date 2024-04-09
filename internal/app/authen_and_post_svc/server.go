package authen_and_post_svc

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/leducanh112/newsfeed/configs"
	"github.com/leducanh112/newsfeed/pkg/types/proto/pb/authen_and_post"
)

type AuthenticateAndPostServer struct {
	authen_and_post.UnimplementedAuthenticateAndPostServer
	db    *gorm.DB
	redis *redis.Client

	log *zap.Logger
}

func NewAuthenticateAndPostService(conf configs.AuthenticateAndPostConfig) (*AuthenticateAndPostServer, error) {
	db, err := gorm.Open(mysql.New(conf.MySQL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("can not connect to db ", err)
		return nil, err
	}
	rd := redis.NewClient(&conf.Redis)
	if rd == nil {
		return nil, fmt.Errorf("can not init redis client")
	}

	log, err := newLogger()
	if err != nil {
		return nil, err
	}
	return &AuthenticateAndPostServer{
		db:    db,
		redis: rd,
		log:   log,
	}, nil
}

func newLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to init zap")
	}
	return logger, nil
}
