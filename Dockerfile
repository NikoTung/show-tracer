FROM golang:1.17-alpine as builder
RUN apk add git openssh-client
WORKDIR /go/app
COPY . .
RUN mkdir -p /root/.ssh \
    && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

RUN go mod tidy
RUN go build -o app

FROM alpine:3.15.0

WORKDIR /go/app
COPY --from=builder /go/app/app /usr/local/bin/app

CMD ["/usr/local/bin/app", "-f", "/go/app/config.json"]
