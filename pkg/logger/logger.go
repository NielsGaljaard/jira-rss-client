package logger

import zap "go.uber.org/zap"

func New(loglevel string) (*zap.Logger, error) {
	var err error
	var log *zap.Logger
	if loglevel == "production" {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		log, err = config.Build()
	} else {
		log, err = zap.NewDevelopment()
	}
	return log, err
}
