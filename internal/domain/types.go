package domain

type Ticket struct {
	Id        string
	Key       string
	Link      string
	Assignee  string
	Priority  string
	Reporter  string
	IssueType string
	Status    string
	Summary   string
}
type Channel struct {
	JQL         string `mapstructure:"JQL"`
	Title       string `mapstructure:"Title"`
	Link        string `mapstructure:"Link"`
	Description string `mapstructure:"Description"`
	FilePath    string `mapstructure:"Filepath"`
}
