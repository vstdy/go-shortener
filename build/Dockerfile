FROM golang:1.17.7-alpine3.15 AS build

WORKDIR /go/src/github.com/vstdy0/go-project
COPY . .
RUN go build -o shortener ./cmd/shortener/main.go

FROM alpine:3.15

RUN apk add --no-cache net-tools

WORKDIR /go/src/github.com/vstdy0/go-project
COPY --from=build /go/src/github.com/vstdy0/go-project/shortener .
COPY ./storage/psql/migrations/*.sql ./storage/psql/migrations/
COPY ./build/entrypoint.sh .

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD ["./shortener"]
