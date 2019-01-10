package trelloclient

import "app/payloads"
import "fmt"
import "log"
import "regexp"
import "strconv"

import "github.com/mauromsl/trello"

const issuePrefixRE string = "^#[0-9]+: "

type trelloClient struct {
	*trello.Client
	ListId     string
	BoardId    string
	issueCards map[int64]string
}

var clientCache = make(map[string]*trelloClient)

func (c trelloClient) NewIssue(issue payloads.GithubIssue) error {
	if _, ok := c.issueCards[issue.Number]; ok {
		log.Println("Card already exists for issue: #", issue.Number)
		return nil
	}
	log.Println("Creating card for issue: #", issue.Number)
	attachment := trello.Attachment{
		URL:  issue.HTMLURL,
		Name: issue.Title,
	}
	card := trello.Card{
		Name:        fmt.Sprintf("#%d: %s", issue.Number, issue.Title),
		Desc:        fmt.Sprintf("url: $s", issue.HTMLURL),
		Pos:         1,
		IDList:      c.ListId,
		Attachments: []*trello.Attachment{&attachment},
	}
	log.Println("Creating Card: ", card)
	err := c.CreateCard(&card, trello.Defaults())
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	c.issueCards[issue.Number] = card.ID
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

func (c trelloClient) MakeCardIndex() (map[int64]string, error) {
	index := make(map[int64]string)
	r, _ := regexp.Compile(issuePrefixRE)
	board, err := c.GetBoard(c.BoardId, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		match := r.FindString(card.Name)
		log.Println("> ", card.Name, match)
		if match != "" {
			id, _ := strconv.ParseInt(match[1:len(match)-2], 10, 64)
			index[id] = card.ID
		}
	}

	return index, err
}

func NewClient(apiKey, oauthToken, boardId, listId string) *trelloClient {
	client := &trelloClient{Client: trello.NewClient(apiKey, oauthToken)}
	client.ListId = listId
	client.BoardId = boardId
	issueCards, err := client.MakeCardIndex()
	if err != nil {
		panic(err)
	}
	client.issueCards = issueCards
	log.Println("Created Trello client with %v issue cards", len(issueCards))
	return client
}

func GetClient(apiKey, oauthToken, boardId, listId string) *trelloClient {
	client, cached := clientCache[apiKey+oauthToken+listId]
	if !cached {
		client = NewClient(apiKey, oauthToken, boardId, listId)
		clientCache[apiKey+oauthToken+listId] = client
	}
	return client
}
