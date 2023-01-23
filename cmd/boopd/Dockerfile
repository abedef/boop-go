FROM golang:1.16

WORKDIR /usr/src/app

# see: https://hub.docker.com/_/golang
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/boopd

EXPOSE 22022

CMD ["app"]
