# 1 step
FROM golang:1.22-alpine AS builder

# . means current repo
COPY . /github.com/AndreiMartynenko/chat-server/source/
WORKDIR /github.com/AndreiMartynenko/chat-server/source/

RUN go mod download
RUN go build -o ./bin/crud_server ./cmd/grpc_server/main.go

# 2 step
FROM alpine:latest
WORKDIR /root/

# . here is copy everything in the root
COPY --from=builder /github.com/AndreiMartynenko/chat-server/source/bin/crud_server .

# run our server
CMD ["./crud_server"]