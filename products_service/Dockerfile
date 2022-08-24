FROM golang:alpine AS builder

ENV CGO_ENABLED=0

RUN mkdir /app
WORKDIR /app

COPY commons/ /commons/

COPY products_service/go.* ./

RUN go mod download

COPY products_service/. .

RUN go build -o ./app cmd/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]