# Overview
This was a for-fun project I did as a counterpart to a websocket based, asynchronous browser event distribution platform.
This implementation uses grpc streams and illustrates how to use the streaming server to do this. Since the gRPC call
for `Subscribe` is called in the gRPC code from a goroutine `Subscribe` has to capture that goroutine and hold the server
instance which can then deliver messages that come in out-of-band to the client. If `Subscribe` returns, the context
is closed.

This method is a bit cleaner and easier that the websocket design but undoubtedly uses more memory but not sure how
much more or if that matters a ton. Sometimes the context closes well after a tab closes and depending on the browser, 
but it eventually does. Websocket version closes instantly on client hangup. In reality this would still need a side channel distribution network like redis or 
maybe RabbitMQ (AMQP).

Since we are using gRPC-web a proxy is required, so an Envoy instance is included but any proxy should work that supports
gRPC-web.

# Install and Run
```
brew install go-grpc protoc-gen-go-grpc protoc-gen-grpc-web
#In different terminals
make run-proxy
make run-server
make run-client
```

# Profiling
The main is set up to profile. Long story short is most of the memory is allocated in the Subscribe function where 
the client channel is allocated. I was able to serve 250K clients in under 1GB. 

```bash
go tool pprof --pdf server/server server/profiles/<profile>.pprof > <profile>.pdf
```
