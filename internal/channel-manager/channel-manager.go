package channelmanager

import (
	"github.com/NielsGaljaard/jira-rss-client/internal/domain"
	rss_writer "github.com/NielsGaljaard/jira-rss-client/internal/rss-writer"
	jira_client "github.com/NielsGaljaard/jira-rss-client/pkg/jira-client"
	"go.uber.org/zap"

	"os"
)

type Config struct {
	Channels *[]domain.Channel
	Writer   *rss_writer.Writer
	Client   *jira_client.Client
}

type ChannelManager struct {
	logger   *zap.Logger
	channels *[]domain.Channel
	writer   *rss_writer.Writer
	client   *jira_client.Client
}

func New(c *Config, logger *zap.Logger) (*ChannelManager, error) {
	return &ChannelManager{
		channels: c.Channels,
		writer:   c.Writer,
		client:   c.Client,
		logger:   logger,
	}, nil
}

func (c *ChannelManager) UpdateChannels() {
	for _, channel := range *c.channels {
		c.logger.Debug("updating", zap.Any("channel", channel))
		tickets, err := c.client.DoJQL(channel.JQL)
		if err != nil {
			c.logger.Error("failed to do JQL fetch for channel", zap.String("channel", channel.Title), zap.Error(err))
		}
		file, err := os.Create(channel.FilePath)
		defer file.Close()
		if err != nil {
			c.logger.Error("failed to open file for channel", zap.String("channel", channel.Title), zap.Error(err))
		}
		c.logger.Info("writing tickets for", zap.Int("nr_tickets", len(*tickets)), zap.String("channel", channel.Title), zap.String("JQL", channel.JQL))
		err = c.writer.WriteToPath(&rss_writer.RssInput{
			Description: channel.Description,
			Link:        channel.Link,
			Title:       channel.Title,
			Tickets:     tickets,
		}, file)
		if err != nil {
			c.logger.Error("error writing tickets for channel", zap.String("channel", channel.Title), zap.Error(err))
		}
	}
}
