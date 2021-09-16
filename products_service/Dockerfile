FROM golang:alpine AS builder

ENV CGO_ENABLED=0

RUN mkdir /app
WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o ./app

FROM scratch

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]