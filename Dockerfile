FROM golang:1.16

WORKDIR /go/src/app
RUN mkdir -p /opt/data

COPY src .
COPY go.mod .
COPY README.md .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["notify_ldn_inbox","db=","/opt/data/ldn_inbox.sqlite"]