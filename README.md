![Intersect Logo](https://raw.githubusercontent.com/mauromsl/Intersect/master/intersect-logo.png)
# Intersect

Intersect allows you to syncronise a github project with a trello board

## Supported features:

 - Creation of trello cards for new open issues
 - Creation of comments on those cards when they are added to their respective comment
 - Integration with the [GitHub power up for Trello]( https://trello.com/power-ups/55a5d916446f517774210004/github)

## Upcoming features
 - Synchronization of labels across Github and Trello
 - Support for syncronitinsing `edit` and `delete` actions on GitHub comments


# Installation Instructions
Intersect is a go project built on Go 1.11 and can be built as a regular go project form `$GOPATH/src/` with `go get && go install && go build`
The application takes the following flags at startup

- `--oauth-token` (Your Trello OAUTH token)
- `--api-key` (Your Trello API Key)
- `--board-id` (The ID of the board where new cards will be created)
- `--list-id` (The ID of the board in which to place the newly created cards)

GNU cMake targets are provided for building and running the project within docker containers:

## Steps
 - Export the previously mentioned flags via `make` variables or environment variables (`$GI_OAUTH_TOKEN`, `$GI_API_KEY`, `$GI_BOARD_ID`, `$GI_LIST_ID`)
 - Run `make image` to build the base docker image
 - Run `make` to compile and run the application (alternatively run `make build` to compile and `make run` to run)
