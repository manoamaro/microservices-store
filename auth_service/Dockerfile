FROM golang:alpine AS builder

ENV CGO_ENABLED=0

RUN mkdir /app
WORKDIR /app

COPY commons/ /commons/

COPY auth_service/go.* ./

RUN go mod download

COPY auth_service/. .

RUN go build -o ./app

FROM scratch

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]