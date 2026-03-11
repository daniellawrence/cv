FROM golang:1.25 AS builder

WORKDIR /workspace

# copy workspace files
COPY go.work ./
COPY go.work.sum ./

# copy modules
COPY backend ./backend
COPY gen ./gen

WORKDIR /workspace/backend/education

RUN go mod download

ARG CONTENT_SHA1=unknown
ARG BUILD_TIMESTAMP=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X 'github.com/daniellawrence/cv/backend/common.ContentSHA1=${CONTENT_SHA1}' \
              -X 'github.com/daniellawrence/cv/backend/common.BuildTimestamp=${BUILD_TIMESTAMP}'" \
    -o app .

FROM scratch
COPY --from=builder /workspace/backend/education/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]