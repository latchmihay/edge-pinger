FROM golang:latest as build

ARG LD_FLAGS

WORKDIR /go/src/github.com/latchmihay/edge-pinger
COPY . .

RUN go version
#RUN go test ./...
RUN GO111MODULE=on CGO_ENABLED=0 go build -ldflags "$LD_FLAGS" -o /go/bin/app .

FROM gcr.io/distroless/static
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
