# build stage
FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY . .

RUN go build -o main main.go

RUN apk update && apk add --no-cache curl tar

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz 


# run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY ./app.env .
COPY db/migration ./migration
COPY start.sh .
COPY wait-for.sh .


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]