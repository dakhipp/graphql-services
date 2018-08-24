FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/auth
COPY vendor ../vendor
COPY auth ./
RUN go build -o /go/bin/app ./cmd/auth/main.go
RUN go build -o /go/bin/migrates ./migrations

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD migrates migrate ; app
