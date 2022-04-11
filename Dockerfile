##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /application_metadata_api_server

COPY . /application_metadata_api_server

RUN go mod download
RUN go build -o  /bin/application_metadata_api_server

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bin/application_metadata_api_server /bin/application_metadata_api_server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bin/application_metadata_api_server"]

