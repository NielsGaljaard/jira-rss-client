package rss_writer

import (
	"io"
	tex "text/template"

	"github.com/NielsGaljaard/jira-rss-client/internal/domain"
	"go.uber.org/zap"
)

type Writer struct {
	logger   *zap.Logger
	template *tex.Template
}

type RssInput struct {
	Title       string
	Description string
	Link        string
	Tickets     *[]domain.Ticket
}

func New(templatePath string, logger *zap.Logger) (*Writer, error) {
	rssTemplate, err := tex.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return &Writer{
		logger:   logger,
		template: rssTemplate,
	}, nil
}

func (w *Writer) WriteToPath(rss *RssInput, writer io.Writer) error {
	err := w.template.Execute(writer, rss)
	if err != nil {
		w.logger.Error("error writing template", zap.Error(err))
	}
	return err
}
