package payloads

const NEW_ISSUE = "opened"
const NEW_ISSUE_COMMENT = "created"

type IssuesEventPayload struct {
	Action     string           `json:"action"`
	Issue      GithubIssue      `json: issue`
	Repository GithubRepository `json: repository`
	Sender     GithubSender     `json: sender`
	Assignee   *GithubAssignee  `json:"assignee"`
}

type IssueCommentPayload struct {
	Action     string           `json:"action"`
	Issue      GithubIssue      `json:"issue"`
	Comment    GithubComment    `json:"comment"`
	Repository GithubRepository `json:"repository"`
	Sender     GithubUser       `json:"sender"`
}
