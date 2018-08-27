FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/migrations
COPY vendor ../vendor
COPY migrations ./
RUN go build -o /go/bin/migrates ./

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD migrates migrate
