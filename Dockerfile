FROM golang
ADD . /go/src/github.com/pagerduty/go-pagerduty
WORKDIR /go/src/github.com/pagerduty/go-pagerduty
RUN go get ./... && go test -v -race -cover ./...
