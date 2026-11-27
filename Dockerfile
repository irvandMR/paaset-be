# ---------- STAGE 1: Builder ----------
FROM golang:1.26-alpine AS builder        
WORKDIR /app
COPY go.mod go.sum ./                   
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o paaset-be ./cmd/paaset 

# ---------- STAGE 2: Runner ----------
FROM alpine:3.20                          
RUN apk add --no-cache ca-certificates tzdata  
WORKDIR /app
COPY --from=builder /app/paaset-be .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./paaset-be"]
