FROM golang:1.16.7-alpine3.14

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/sbpann/go-docker-multi-stage-build-graceful-shutdown-example

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:3.14
# optional for gin-gonic
ENV GIN_MODE=release
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/src/github.com/sbpann/go-docker-multi-stage-build-graceful-shutdown-example.
ENTRYPOINT ["./app"]
