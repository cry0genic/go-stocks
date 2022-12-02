FROM golang:1.16-alpine as builder
RUN apk --update upgrade
RUN apk add --no-cache git make musl-dev gcc sqlite libc6-compat
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod download
RUN GOOS=linux CGO_ENABLED=1 go build -a -o stocks .

FROM alpine:latest
RUN apk --update upgrade
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /app/stocks .

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
CMD ["/stocks"]