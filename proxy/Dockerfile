FROM envoyproxy/envoy-alpine:v1.19-latest
COPY config.yaml /etc/config.yaml
CMD ["/usr/local/bin/envoy", "-c", "/etc/config.yaml"]
