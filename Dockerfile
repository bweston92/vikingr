FROM golang:1.9-alpine
RUN apk update && apk add git
COPY . /go/src/github.com/bweston92/vikingr
WORKDIR /go/src/github.com/bweston92/vikingr
RUN go-wrapper download
RUN go build -x -o /tmp/vikingr

FROM alpine:3.6
COPY --from=0 /tmp/vikingr /usr/bin/vikingr
COPY docker-entrypoint.sh /usr/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/bin/docker-entrypoint.sh"]
