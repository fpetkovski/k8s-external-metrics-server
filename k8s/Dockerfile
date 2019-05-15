FROM golang:1.11 as build
COPY . /build
RUN mkdir /out && \
    cd /build && \
    CGO_ENABLED=0 GOARCH=amd64 go build -mod vendor -o /out/adapter

FROM scratch
COPY --from=build /out/adapter /
ENTRYPOINT ["/adapter", "--logtostderr=true"]
