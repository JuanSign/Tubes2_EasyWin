# Build stage
FROM golang:1.24@sha256:39d9e7d9c5d9c9e4baf0d8fff579f06d5032c0f4425cdec9e86732e8e4e374dc AS builder

WORKDIR /app

COPY backend/go.mod ./
RUN go mod download

COPY backend/. .

RUN go build -o app

# Final stage
FROM gcr.io/distroless/base-debian11:latest@sha256:ac69aa622ea5dcbca0803ca877d47d069f51bd4282d5c96977e0390d7d256455

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
