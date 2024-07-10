# Chat Service

Chat Service includes 3 docker containers:
- chat-server - controls chats.
- postgres - permanent storage of data.
- migrator - performs migration in database using goose package.

## Deploy

Make sure docker network `service-net` is in place for microservices communication. If none exists, then create network:
```
# make docker-net
```

To deploy Chat Service:
```
# make docker-deploy ENV=<environment>
```
*ENV is used then as a config name. Possible ENV values are now `stage` and `prod` as these configs are now in the repository.*

To stop Chat Service:
```
# make docker-stop ENV=<environment>
```

# Chat Client

Simple cli utility for communicating with each other through chats can be built with:
```
make build-chat-client
```

Usage examples:

```
bin/chat-client login --address=localhost:50000 --cert=tls/auth.pem
bin/chat-client logout

bin/chat-client create chat --address=localhost:50001 --cert=tls/chat.pem --users=Alice,Bob,Carol
bin/chat-client delete chat --address=localhost:50001 --cert=tls/chat.pem --id=45d8c8ce-818d-40bc-8e6c-6afd28cf34b4
bin/chat-client connect chat --address=localhost:50001 --cert=tls/chat.pem --id=45d8c8ce-818d-40bc-8e6c-6afd28cf34b4
```