package main

import "app/payloads"
import "app/intersect"
import "net/http"
import "log"
import "os"
import "encoding/json"
import "flag"
import "os/signal"
import "github.com/gorilla/mux"

const ISSUE_EVENT string = "issues"
const ISSUE_COMMENT_EVENT string = "issue_comment"

var apiKey string
var oauthToken string
var boardId string
var listId string

func init() {
	flag.StringVar(&apiKey, "api-key", "", "Trello API Key")
	flag.StringVar(&oauthToken, "oauth-token", "", "Trello OAUTH Token")
	flag.StringVar(&boardId, "board-id", "", "Identifier of the Trello Board to be synced with github")
	flag.StringVar(&listId, "list-id", "", "Identifier of the Trello board list in which new cards are created")

}

func validateFlags() {
	flag.Parse()
	if flag.NFlag() != 4 {
		log.Println("Some flags are not set")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func home(response http.ResponseWriter, request *http.Request) {
	text := "Connected"
	response.Write([]byte(text))
}

func github(response http.ResponseWriter, request *http.Request) {
	intersectClient := intersect.GetClient(apiKey, oauthToken, boardId, listId)
	eventType := request.Header.Get("X-Github-Event")
	log.Println("Incoming Request: ", request)
	decoder := json.NewDecoder(request.Body)
	switch eventType {
	case ISSUE_EVENT:
		var event payloads.IssuesEventPayload
		err := decoder.Decode(&event)
		if err != nil {
			syntax := err.(*json.SyntaxError)
			log.Fatalf("Error decoding JSON: ", syntax, syntax.Offset)
		}
		intersectClient.HandleIssue(event)
	case ISSUE_COMMENT_EVENT:
		var event payloads.IssueCommentPayload
		err := decoder.Decode(&event)
		if err != nil {
			syntax := err.(*json.SyntaxError)
			log.Fatalf("Error decoding JSON: ", syntax, syntax.Offset)
		}
		intersectClient.HandleIssueComment(event)
	default:
		log.Printf("Ignorig event of type: %s", eventType)
	}
	return
}

func main() {
	validateFlags()

	router := mux.NewRouter()
	router.HandleFunc("/github", github).
		Methods("POST")
	router.HandleFunc("/", home).
		Methods("GET")

	app := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go run(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Exiting")
	os.Exit(0)
}

func run(app *http.Server) {
	log.Println("Listening to requests on port 8080")
	if err := app.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
