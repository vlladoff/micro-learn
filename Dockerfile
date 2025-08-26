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

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /build/app /app

ENV TZ=UTC

EXPOSE 8080

ENTRYPOINT ["/app"]