APP_NAME=intersect

RESTART_POLICY ?= no
REMOVE ?= --rm
MODE ?= -it
all:
	make build && make run
image:
	docker build -t $(APP_NAME) --no-cache .
	@notify-send 'Image Built' 2>/dev/null | true
build:
	docker run --rm -v "$(PWD)/$(APP_NAME)":/usr/src/myapp -w /usr/src/myapp $(APP_NAME) go build -v
	@notify-send 'Gitersect Compiled' 2>/dev/null | true
run:
	docker run $(MODE) -p 8080:8080 $(REMOVE) --name $(APP_NAME) -v $(PWD)/$(APP_NAME):/go/src/app \
		--restart $(RESTART_POLICY) \
		--entrypoint /go/src/app/myapp $(APP_NAME) \
			--oauth-token $(GI_OAUTH_TOKEN) \
			--api-key $(GI_API_KEY) \
			--board-id $(GI_BOARD_ID) \
			--list-id $(GI_LIST_ID)

daemon:
	bash -c "make run RESTART_POLICY=always REMOVE='' MODE=-d"

shell:
	docker run -it -v $(PWD)/$(APP_NAME):/go/src/app --rm $(APP_NAME) /bin/bash
