package logger

import zap "go.uber.org/zap"

func New(loglevel string) (*zap.Logger, error) {
	var err error
	var log *zap.Logger
	if loglevel == "production" {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}
	return log, err
}
