package intersect

import "app/payloads"
import "log"
import "regexp"
import "strconv"

import "github.com/mauromsl/trello"

const issuePrefixRE string = "^#[0-9]+: "

type intersectClient struct {
	Trello     *trello.Client
	ListId     string
	BoardId    string
	issueCards map[int64]string
}

var clientCache = make(map[string]*intersectClient)

func (c intersectClient) HandleAction(event payloads.IssuesEventPayload) error {
	var err error
	issue := event.Issue
	switch event.Action {
	case payloads.NEW_ISSUE:
		err = c.NewIssue(issue)
	default:
		log.Printf("Ignoring action: %s", event.Action)
	}
	return err
}

func (c intersectClient) MakeCardIndex() (map[int64]string, error) {
	index := make(map[int64]string)
	r, _ := regexp.Compile(issuePrefixRE)
	board, err := c.Trello.GetBoard(c.BoardId, trello.Defaults())
	if err != nil {
		return nil, err
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		match := r.FindString(card.Name)
		if match != "" {
			id, _ := strconv.ParseInt(match[1:len(match)-2], 10, 64)
			index[id] = card.ID
		}
	}

	return index, err
}

func NewClient(apiKey, oauthToken, boardId, listId string) *intersectClient {
	client := &intersectClient{
		Trello:  trello.NewClient(apiKey, oauthToken),
		ListId:  listId,
		BoardId: boardId,
	}
	issueCards, err := client.MakeCardIndex()
	if err != nil {
		panic(err)
	}
	client.issueCards = issueCards
	log.Printf("Created Trello client for board %s with %d issue cards", boardId, len(issueCards))
	return client
}

func GetClient(apiKey, oauthToken, boardId, listId string) *intersectClient {
	client, cached := clientCache[apiKey+oauthToken+listId]
	if !cached {
		client = NewClient(apiKey, oauthToken, boardId, listId)
		clientCache[apiKey+oauthToken+listId] = client
	}
	return client
}
