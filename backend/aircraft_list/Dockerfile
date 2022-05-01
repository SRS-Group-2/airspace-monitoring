
##
## Build
##
FROM golang:latest AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -a -o /app.exe

##
## Deploy
##
FROM gcr.io/distroless/base:latest

WORKDIR /

COPY --from=build /app.exe /app.exe

ENV GIN_MODE release
ENV PORT 9032

EXPOSE 9032

USER nonroot:nonroot

ENTRYPOINT ["/app.exe"]

## TODO: Need to copy the google credentials file into the release container or it won't be read.
## Also container version won't work with pubsub emulator since they are running on different machines.

## docker build . --tag aircraft_list
## docker run --publish 8080:9032 --env-file ./env.list aircraft_list