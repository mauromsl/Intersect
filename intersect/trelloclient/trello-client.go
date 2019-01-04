package trelloclient

import "app/payloads"
import "fmt"
import "log"

import "github.com/mauromsl/trello"

type trelloClient struct {
	*trello.Client
	ListId string
}

func (c trelloClient) NewIssue(issue payloads.GithubIssue) error {
	log.Println("Creating card for issue: #", issue.Number)
	attachment := trello.Attachment{
		URL: issue.HTMLURL,
		Name: issue.Title,
	}
	card := trello.Card{
		Name:     fmt.Sprintf("#%d: %s", issue.Number, issue.Title),
		Desc:     fmt.Sprintf("url: $s", issue.HTMLURL),
		Pos:      1,
		IDList:   c.ListId,
		Attachments: []*trello.Attachment{&attachment},
	}
	log.Println("Creating Card: ", card)
	err := c.CreateCard(&card, trello.Defaults())
	if err != nil {
		log.Println("Error: ", err)
	}
	log.Println("Attaching Github issue to card: ", card)
	err = card.AddURLAttachment(&attachment)
	if err != nil {
		log.Println("Error: ", err)
	}
	return err
}

func (c trelloClient) HandleAction(event payloads.IssuesEventPayload) error {
	var err error
	issue := event.Issue
	switch event.Action {
	case payloads.NEW_ISSUE:
		err = c.NewIssue(issue)
	default:
		log.Println("Ignoring action: ", event.Action)
	}
	return err
}

func NewClient(apiKey string, oauthToken string, listId string) *trelloClient {
	client := &trelloClient{Client: trello.NewClient(apiKey, oauthToken)}
	client.ListId = listId
	return client
}
