FROM golang:1.18 as builder

# first (build) stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o owl

# final (target) stage

FROM alpine:3.16.0
COPY --from=builder /app/owl /
COPY db db
CMD ["/owl", "server", "--port", "8080"]
EXPOSE 8080
