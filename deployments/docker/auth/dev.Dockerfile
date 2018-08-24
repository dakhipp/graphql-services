FROM golang:1.10.2-alpine3.7
RUN apk --no-cache add git gcc g++ make ca-certificates
RUN go get github.com/dakhipp/go-pg-migrations
RUN go get github.com/oxequa/realize
WORKDIR /go/src/github.com/dakhipp/graphql-services/auth
COPY vendor ../vendor
COPY auth ./
CMD go run ./migrations/*.go migrate ; realize start
