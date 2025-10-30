# ---------- Build Stage ----------
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# install dependencies
RUN apk add --no-cache git make protoc protobuf-dev

# install go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# make sure go bin is in path
ENV PATH="/go/bin:${PATH}"

# copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# generate protobuf files
RUN make build-protobuf

# build binaries
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/migrate cmd/migrate/main.go

# ---------- Runtime Stage ----------
FROM alpine:latest

# install certs for HTTP requests
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# copy binaries, configs, and migrations
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# expose application port
EXPOSE 9000

# run the service
ENTRYPOINT ["./bin/app"]