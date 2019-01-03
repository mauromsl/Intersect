APP_NAME=intersect

all:
	make build && make run
image:
	docker build -t $(APP_NAME) --no-cache .
	@notify-send 'Image Built'
build:
	docker run --rm -v "$(PWD)/$(APP_NAME)":/usr/src/myapp -w /usr/src/myapp $(APP_NAME) go build -v
	@notify-send 'Gitersect Compiled' 2>/dev/null
run:
	bash -c "source $(PWD)/.env"
	docker run -it -p 8080:8080 --rm --name $(APP_NAME) -v $(PWD)/$(APP_NAME):/go/src/app --entrypoint /go/src/app/myapp $(APP_NAME) \
		--oauth-token $(GI_OAUTH_TOKEN) \
		--api-key $(GI_API_KEY) \
		--board-id $(GI_BOARD_ID) \
		--list-id $(GI_LIST_ID)
shell:
	docker run -it -v $(PWD)/$(APP_NAME):/go/src/app --rm $(APP_NAME) /bin/bash
