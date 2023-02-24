FROM golang:alpine AS builder

ARG BUILD_FOLDER
ENV CGO_ENABLED=0

RUN apk update && apk add git

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./app ${BUILD_FOLDER}/cmd/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]