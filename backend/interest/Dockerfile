FROM golang:1.23 AS builder

WORKDIR /workspace

# copy workspace files
COPY go.work ./
COPY go.work.sum ./

# copy modules
COPY backend ./backend
COPY gen ./gen

WORKDIR /workspace/backend/interest

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM scratch
COPY --from=builder /workspace/backend/interest/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]