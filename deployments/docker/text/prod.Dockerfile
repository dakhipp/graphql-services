FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/text
COPY vendor ../vendor
COPY text ./
RUN go build -o /go/bin/app ./cmd/text/main.go

FROM golang:1.10.2-alpine3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
COPY text/.env.dev ./.env.dev
CMD ["app"]
