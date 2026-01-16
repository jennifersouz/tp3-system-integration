FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o xml-service main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/xml-service /xml-service
EXPOSE 8081
ENTRYPOINT ["/xml-service"]
