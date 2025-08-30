FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a -installsuffix cgo \
  -o app ./cmd/smpl-api-oapi/

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /build/app /app

ENV TZ=UTC

EXPOSE 8080

ENTRYPOINT ["/app"]