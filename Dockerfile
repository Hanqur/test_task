FROM golang:alpine3.16 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go ./

ENV CGO_ENABLED=0

RUN go build -o /application

FROM scratch

COPY --from=build /application /application

ENV APP_LOG_LEVEL="info" \
    APP_DB_FILE_NAMES="db.json,db.csv"

CMD [ "/application" ]

