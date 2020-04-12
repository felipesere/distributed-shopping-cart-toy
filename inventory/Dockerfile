# Start by building the application.
FROM golang:alpine as build

# All these steps will be cached
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

# Now we add the source
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app cmd/server/server.go

# Now copy it into our base image.
FROM scratch
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
