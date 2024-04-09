package service

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/leducanh112/newsfeed/configs"
	authen_and_post2 "github.com/leducanh112/newsfeed/pkg/client/authen_and_post"
	"github.com/leducanh112/newsfeed/pkg/types/proto/pb/authen_and_post"
)

type WebService struct {
	authenticateAndPostClient authen_and_post.AuthenticateAndPostClient

	log *zap.Logger
}

func NewWebService(conf configs.WebConfig) (*WebService, error) {
	aapClient, err := authen_and_post2.NewClient(conf.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, fmt.Errorf("failed to init AAP client: %s", err)
	}

	log, err := newLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to init log: %s", err)
	}
	return &WebService{
		authenticateAndPostClient: aapClient,
		log:                       log,
	}, nil
}

func newLogger() (*zap.Logger, error) {
	// f, err := os.OpenFile("./logs/webapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create log file")
	// }
	// f.Close()
	//
	// cfg := zap.Config{
	// 	Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
	// 	Encoding:         "console",
	// 	OutputPaths:      []string{"stdout", "./logs/webapp.log"},
	// 	ErrorOutputPaths: []string{"stderr", "./logs/webapp.log"},
	// }
	// logger := zap.Must(cfg.Build())
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to init zap")
	}
	return logger, nil
}
