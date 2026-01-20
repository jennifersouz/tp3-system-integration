FROM golang:1.22 AS build
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends protobuf-compiler ca-certificates && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# gera os pb dentro de /app/pb
RUN mkdir -p pb && \
    protoc -I proto \
      --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      proto/bi_grpc.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /xml-service .

FROM gcr.io/distroless/base-debian12:latest
COPY --from=build /xml-service /xml-service
EXPOSE 8081
EXPOSE 50051
EXPOSE 50052
ENTRYPOINT ["/xml-service"]
