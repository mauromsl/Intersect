FROM golang:1.11

WORKDIR /go/src/app
COPY ./intersect .

EXPOSE 8080

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
