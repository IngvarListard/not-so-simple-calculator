FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

COPY . /go/src/calculator
WORKDIR /go/src/calculator/cmd/calculator

RUN apk add build-base
RUN go mod vendor
RUN go test /go/src/calculator/...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o calculator

FROM alpine:3

RUN apk --update upgrade && \
    apk add sqlite

COPY --from=builder /go/src/calculator/cmd/calculator/calculator /go/bin/calculator
RUN chmod +x /go/bin/calculator
RUN mkdir -p /srv/data/database

CMD ["./go/bin/calculator"]