FROM golang:1.10.2-alpine3.7
RUN apk --no-cache add git gcc g++ make ca-certificates
RUN go get github.com/oxequa/realize
WORKDIR /go/src/github.com/dakhipp/graphql-services/graphql
COPY vendor ../vendor
COPY auth ../auth
COPY graphql ./
CMD ["realize", "start"]
