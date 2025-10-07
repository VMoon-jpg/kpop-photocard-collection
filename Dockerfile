FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/sample-cards.jsonl ./

# Create cards.jsonl from sample if it doesn't exist
RUN if [ ! -f cards.jsonl ]; then cp sample-cards.jsonl cards.jsonl; fi

EXPOSE 8080

CMD ["./main"]