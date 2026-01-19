FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod init xmlservice 2>/dev/null || true && \
    go mod tidy 2>&1 || true && \
    CGO_ENABLED=0 GOOS=linux go build -o xml-service main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/xml-service /xml-service
EXPOSE 8081 50051 50052
EXPOSE 8081
ENTRYPOINT ["/xml-service"]
