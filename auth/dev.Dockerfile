FROM golang:1.10.2-alpine3.7
RUN apk --no-cache add git gcc g++ make ca-certificates
RUN go get github.com/codegangsta/gin
WORKDIR /go/src/github.com/dakhipp/graphql-services/auth
COPY vendor ../vendor
COPY auth ./
CMD ["gin", "run", "./cmd/auth/main.go"]
