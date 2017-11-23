FROM golang:1.9 AS builder
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep
WORKDIR /go/src/github.com/jonahgeorge/weatherglass/
COPY . /go/src/github.com/jonahgeorge/weatherglass/
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 go build -a -installsuffix cgo

FROM alpine
CMD ["/weatherglass"]
COPY --from=builder /go/src/github.com/jonahgeorge/weatherglass/public /public
COPY --from=builder /go/src/github.com/jonahgeorge/weatherglass/templates /templates
COPY --from=builder /go/src/github.com/jonahgeorge/weatherglass/weatherglass /weatherglass
