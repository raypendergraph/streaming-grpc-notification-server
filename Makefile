GO=go
NPM=npm
ROOT_DIR := $(shell pwd)
JS_OUT=$(ROOT_DIR)/client/src/generated
GO_OUT=$(ROOT_DIR)/server/generated

.PHONY: server-proto client-proto run-server run-client

server/server: server-proto
	cd $(ROOT_DIR)/server \
		&& $(GO) build \
		&& chmod 744 server

run-server: server/server
	cd $(ROOT_DIR)/server && ./server

server-proto:
	-mkdir -p $(GO_OUT)
	protoc \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT)\
		proto/notifications.proto

client-proto:
	-mkdir $(JS_OUT)
	protoc \
		--js_out=import_style=commonjs:$(JS_OUT)\
    	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(JS_OUT) \
		proto/notifications.proto

run-client: client-proto
	cd $(ROOT_DIR)/client && $(NPM) i && $(NPM) start	

run-proxy:
	cd $(ROOT_DIR)/proxy \
		&& docker build -t grpc-proxy . \
		&& docker run -p 10000:10000 -p 9901:9901 -it grpc-proxy
clean:
	rm -rf $(GO_OUT)
	rm -rf $(JS_OUT)
	rm server/server
