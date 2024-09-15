FROM golang:1.23.1-bookworm AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify
RUN go install github.com/a-h/templ/cmd/templ@v0.2.778

COPY . .
RUN go generate ./... && CGO_ENABLED=0 go build -v -o /usr/local/bin/app ./cmd

FROM alpine:3.20.3
EXPOSE 8080

COPY --from=build /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
