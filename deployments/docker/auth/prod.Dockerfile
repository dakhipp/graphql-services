FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/auth
COPY vendor ../vendor
COPY email ../email
COPY text ../text
COPY auth ./
RUN go build -o /go/bin/app ./cmd/auth/main.go

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD ["app"]
