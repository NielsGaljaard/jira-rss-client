package main

import (
	"fmt"
	"os"

	channelmanager "github.com/NielsGaljaard/jira-rss-client/internal/channel-manager"
	config_reader "github.com/NielsGaljaard/jira-rss-client/internal/config-reader"
	rss_writer "github.com/NielsGaljaard/jira-rss-client/internal/rss-writer"
	jira_client "github.com/NielsGaljaard/jira-rss-client/pkg/jira-client"
	logger2 "github.com/NielsGaljaard/jira-rss-client/pkg/logger"
	"go.uber.org/zap"
)

const path = "/.config/jrrss/"

func main() {
	log, err := logger2.New("debug")
	if err != nil {
		fmt.Println("FATAL ERROR")
		os.Exit(1)
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cant find home folder for user", zap.Error(err))
	}
	config, err := config_reader.LoadConfig(dirname + path)
	if err != nil {
		log.Fatal("cant open config file", zap.String("fpath", dirname+path), zap.Error(err))
	}
	log, err = logger2.New(config.LogLevel)
	if err != nil {
		fmt.Println("FATAL ERROR")
		os.Exit(1)
	}
	jiraConfig := &jira_client.Config{
		AppUserName: config.AppUser,
		AppPassword: config.AppPassword,
		AppBaseURL:  config.AppUrl,
	}
	client, err := jira_client.New(jiraConfig, log)
	if err != nil {
		log.Fatal("coulnd't retrieve tickets ")
		return
	}
	writer, err := rss_writer.New(config.TemplateLocation, log)
	if err != nil {
		log.Fatal("failed to initialize writer")
	}
	manager, err := channelmanager.New(&channelmanager.Config{
		Writer:   writer,
		Client:   client,
		Channels: &config.Channel,
	}, log)
	if err != nil {
		log.Fatal("failed to initialize channel manager")
	}
	manager.UpdateChannels()
}
