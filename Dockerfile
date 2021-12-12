FROM golang:1.16-alpine as builder

WORKDIR /go/app
COPY . .
RUN go build -o app

FROM golang:1.16-alpine

WORKDIR /go/app
COPY --from=builder /go/app/app /usr/local/bin/app

CMD ["/usr/local/bin/app", "-f", "/go/app/config.json"]