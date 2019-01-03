package trello

import "app/payloads"
import "fmt"
import "log"

import trello_ "github.com/adlio/trello"

//Trello Const
const RequestURL = "https://trello.com/1/OAuthGetRequestToken"
const AccessURL = "https://trello.com/1/OAuthGetAccessToken"
const AuthorizeURL = "https://trello.com/1/OAuthAuthorizeToken"
const AppName = "Intersect"

const BOARD_ID = "5be17c6f3a4edf847029b4a1"
const LIST_ID = "5be17c8d28cb843877bf7313"


type trelloClient struct{
	*trello_.Client
	ListId string
}

func (c trelloClient) NewIssue(issue payloads.GithubIssue) (error){
	card := trello_.Card{
		Name: fmt.Sprintf("#%d: %s", issue.Number, issue.Title),
		Desc: fmt.Sprintf("url: $s", issue.HTMLURL),
		Pos:1,
		IDList: c.ListId,
		IDLabels: []string{"5c0ea689146be2125142b36e"},
	}
	log.Println("Creating Card:", card)
	err := c.CreateCard(&card, trello_.Defaults())
	if err != nil {
		log.Println("Error: ", err)
	}
	return err
}

func (c trelloClient) HandleAction(event payloads.IssuesEventPayload) (error){
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

func NewClient(apiKey string, oauthToken string, listId string)(*trelloClient) {
	client := &trelloClient{Client: trello_.NewClient(apiKey, oauthToken)}
	client.ListId = listId
	return client
}

