# Adopted from jeremyhuiskamp/golang-docker-scratch
FROM golang:alpine as golang
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM scratch
ENV GIN_MODE=release
ENV BINGO_REDIS_ADDR=redis:6379
ENV BINGO_REDIS_DB=4
COPY --from=golang /go/bin/BingoServer /BingoServer
ENTRYPOINT ["/BingoServer"]
