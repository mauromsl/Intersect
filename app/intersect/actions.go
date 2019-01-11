package intersect

import "log"
import "fmt"
import "app/payloads"
import "github.com/mauromsl/trello"

func (c intersectClient) NewIssue(issue payloads.GithubIssue) error {
	if _, ok := c.issueCards[issue.Number]; ok {
		log.Printf("Card already exists for issue: #%d", issue.Number)
		return nil
	}
	log.Println("Creating card for issue: #", issue.Number)
	attachment := trello.Attachment{
		URL:  issue.HTMLURL,
		Name: issue.Title,
	}
	card := trello.Card{
		Name:        fmt.Sprintf("#%d: %s", issue.Number, issue.Title),
		Desc:        fmt.Sprintf("url: %s", issue.HTMLURL),
		Pos:         1,
		IDList:      c.ListId,
		Attachments: []*trello.Attachment{&attachment},
	}
	log.Printf("Creating Card: %s", card.Name)
	err := c.Trello.CreateCard(&card, trello.Defaults())
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	c.issueCards[issue.Number] = card.ID
	log.Printf("Attaching Github issue #%d to card '%s' ", issue.Number, card.Name)
	err = card.AddURLAttachment(&attachment)
	if err != nil {
		log.Println("Error: ", err)
	}
	return err
}
