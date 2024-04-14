FROM golang:1.17-alpine as builder

LABEL maintainer="arkaprovob <apb@live.in>"

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

COPY config.json .

COPY api-specs.yaml .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/api-specs.yaml .

COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./main"]