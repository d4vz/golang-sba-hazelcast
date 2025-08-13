FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /bin/auction ./cmd/auction

FROM gcr.io/distroless/base-debian12

USER nonroot:nonroot
WORKDIR /
COPY --from=builder /bin/auction /auction

ENV APP_HTTP_ADDR=:8080
EXPOSE 8080

ENTRYPOINT ["/auction"]


