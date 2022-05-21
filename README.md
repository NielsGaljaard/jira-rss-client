# Jira Rss Client 

## setup
create a config file under the home config folder and create a storage folder for your rss. And copy over the template (or use your own template)
```sh
mkdir ~/.config/jrrss
touch ~/.config/jrrss/config.yaml
mkdir ~/jrrss
mkdir ~/jrrss/rss/
cp ./internal/rss-writer/rss.gotext ~/jrrss/template.rss
```

## sample configuration
```yaml
AppPassword: <Your app Password>
AppUrl: <Your Jira API URL>
AppUser:  <Your Jira User>
TicketsPerJQL:  <the max number of tickets to be returned per channel>
TemplateLocation: /home/<your user>/jrrss/template.rss
LogLevel: production
Channels:
- JQL: "assignee=currentUser()"
  Title: mytickets
  Link: https://quandago.atlassian.net/issues/?jql=assignee%3DcurrentUser()
  Description: MyTickets
  Filepath: /home/<your user>/jrrss/rss/mytickets.xml
```
## PR and updates
Feel free to contribute, expect little in terms of updates from me

## Product Roadmap
* integrate support for more Atlassian services
* add caching and more complex ticket storage with sqllite
* add support for linked issues
* add support for redering custom fields
* rewrite to rust

