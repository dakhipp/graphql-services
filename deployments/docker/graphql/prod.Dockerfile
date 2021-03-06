FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/graphql
COPY vendor ../vendor
COPY auth ../auth
COPY email ../email
COPY text ../text
COPY graphql ./
RUN go build -o /go/bin/app

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD ["app"]
