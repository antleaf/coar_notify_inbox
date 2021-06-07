FROM golang:1.16

WORKDIR /go/src/app
RUN mkdir -p /opt/data

COPY src .
COPY go.mod .
COPY README.md .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["notify_ldn_inbox","-db","/opt/data/ldn_inbox.sqlite"]
#CMD ["sh", "-c", "notify_ldn_inbox -db=$NOTIFY_LDN_INBOX_DB_PATH -host=$NOTIFY_LDN_INBOX_HOST -port=$NOTIFY_LDN_INBOX_PORT -debug=$NOTIFY_LDN_INBOX_DEBUG"]