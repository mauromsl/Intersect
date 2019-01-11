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
		Desc:        issue.Body,
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

func (c intersectClient) NewComment(comment payloads.GithubComment, issue payloads.GithubIssue) error {
	var err error
	if cardId, ok := c.issueCards[issue.Number]; ok {
		card, err := c.Trello.GetCard(cardId, trello.Defaults())
		if err == nil {
			body := fmt.Sprintf(
				TRELLO_COMMENT_TMPL,
				comment.ID,
				comment.User.Login,
				comment.CreatedAt,
				comment.Body,
			)
			log.Printf(
				"Adding github comment on #%d to trello card: %s",
				issue.Number, card.Name,
			)
			_, err = card.AddComment(body, trello.Defaults())
		}
	} else {
		log.Printf("Received comment for issue that has no card: #%d", issue.Number)
		c.NewIssue(issue)
		c.NewComment(comment, issue)
	}
	return err
}
