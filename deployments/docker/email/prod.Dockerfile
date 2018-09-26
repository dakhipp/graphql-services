FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/dakhipp/graphql-services/email
COPY vendor ../vendor
COPY email ./
RUN go build -o /go/bin/app ./cmd/email/main.go

FROM golang:1.10.2-alpine3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
COPY email/.env.dev ./.env.dev
CMD ["app"]
