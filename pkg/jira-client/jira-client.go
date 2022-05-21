package jira_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/NielsGaljaard/jira-rss-client/internal/domain"
	"go.uber.org/zap"
)

const apiPath = "/rest/api/2/search"

type Config struct {
	TicketsPerJQL uint32
	AppUserName   string
	AppPassword   string
	AppBaseURL    string
}
type Client struct {
	config *Config
	logger *zap.Logger
}

type TicketRequest struct {
	JQL        string   `json:"jql"`
	StartAt    uint32   `json:"startAt"`
	MaxResults uint32   `json:"maxResults"`
	Fields     []string `json:"fields"`
}

type apiResponse struct {
	Tickets []ticketResponse `json:"issues"`
}
type ticketResponse struct {
	Id     string       `json:"id"`
	Key    string       `json:"key"`
	Link   string       `json:"self"`
	Fields ticketFields `json:"fields"`
}
type ticketFields struct {
	Summary   string          `json:"summary"`
	Assignee  ticketAssignee  `json:"assignee"`
	Priority  ticketPriority  `json:"priority"`
	Reporter  ticketReporter  `json:"reporter"`
	IssueType ticketIssueType `json:"issuetype"`
	Status    ticketStatus    `json:"status"`
}
type ticketStatus struct {
	Name string `json:"name"`
}
type ticketAssignee struct {
	DisplayName string `json:"displayName"`
}

type ticketPriority struct {
	Name string `json:"name"`
}

type ticketReporter struct {
	DisplayName string `json:"displayName"`
}
type ticketIssueType struct {
	Name string `json:"name"`
}

func New(config *Config, logger *zap.Logger) (*Client, error) {
	return &Client{
		config: config,
		logger: logger,
	}, nil
}

func (c *Client) DoJQL(JQL string) (*[]domain.Ticket, error) {
	bytes, err := c.doJQLRequest(JQL)
	if err != nil {
		return nil, err
	}
	responses, err := c.parseFromJson(bytes)
	if err != nil {
		return nil, err
	}
	tickets := c.createTicketsFromAPIResponse(responses)
	return tickets, nil
}
func (c *Client) parseFromJson(responseBody []byte) (*[]ticketResponse, error) {
	var response apiResponse
	err := json.Unmarshal(responseBody, &response)
	if err != nil {
		c.logger.Error("error unmarshalling JSON", zap.Error(err))
		c.logger.Info("failed response body", zap.String("body", string(responseBody)))
		return nil, err
	}
	c.logger.Debug("request got", zap.Int("nr_tickets", len(response.Tickets)))
	return &response.Tickets, nil
}
func (c *Client) doJQLRequest(JQL string) ([]byte, error) {
	apiUrl, err := url.Parse(c.config.AppBaseURL + apiPath)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("doing request for", zap.String("url", apiUrl.String()))
	requestBody := &TicketRequest{
		JQL:        JQL,
		StartAt:    0,
		MaxResults: c.config.TicketsPerJQL,
		Fields:     []string{"*navigable"},
	}
	requestByes, err := json.Marshal(requestBody)
	if err != nil {
		c.logger.Error("failed to marshall request body")
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl.String(), bytes.NewReader(requestByes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		c.logger.Error("cant create request", zap.Error(err))
		return nil, err
	}
	req.SetBasicAuth(c.config.AppUserName, c.config.AppPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Error("Api error with body", zap.String("body", string(body)))
			return nil, errors.New("API ERROR")
		}
		c.logger.Error("API error", zap.Error(err), zap.Int("statusCode", resp.StatusCode))
		return nil, errors.New("API ERROR")
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}
func (c *Client) createTicketsFromAPIResponse(t *[]ticketResponse) *[]domain.Ticket {
	tickets := make([]domain.Ticket, len(*t))
	for i, ticket := range *t {
		tickets[i].Id = ticket.Id
		tickets[i].Key = ticket.Key
		tickets[i].Link = ticket.Link
		tickets[i].Assignee = ticket.Fields.Assignee.DisplayName
		tickets[i].Reporter = ticket.Fields.Reporter.DisplayName
		tickets[i].Priority = ticket.Fields.Priority.Name
		tickets[i].IssueType = ticket.Fields.IssueType.Name
		tickets[i].Status = ticket.Fields.Status.Name
		tickets[i].Summary = ticket.Fields.Summary
	}
	return &tickets
}
