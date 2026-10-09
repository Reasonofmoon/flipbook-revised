FROM golang:1.25-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o bin/flipbook .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    libreoffice-impress poppler-utils ca-certificates \
    # CJK fonts so Korean slides don't render as tofu boxes
    fonts-noto-cjk fonts-nanum && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/bin/flipbook .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./flipbook"]
