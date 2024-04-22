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

# Create the appuser
RUN adduser -D -H appuser

WORKDIR /home/appuser

COPY --from=builder /app/main .
COPY --from=builder /app/api-specs.yaml .
COPY --from=builder /app/config.json .

# Change the ownership to appuser
RUN chown -R appuser:appuser /home/appuser

RUN mkdir /home/appuser/store
RUN mkdir /home/appuser/store/files
RUN mkdir /home/appuser/store/record

RUN chown -R appuser:appuser /home/appuser/store
RUN chown -R appuser:appuser /home/appuser/store/files
RUN chown -R appuser:appuser /home/appuser/store/record

RUN chmod -R 777 /home/appuser/store
RUN chmod -R 777 /home/appuser/store/files
RUN chmod -R 777 /home/appuser/store/record

# Set permissions
RUN chmod 755 main config.json api-specs.yaml

EXPOSE 8080

USER appuser
CMD ["./main"]
