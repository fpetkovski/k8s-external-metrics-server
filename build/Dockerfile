FROM golang:1.11 as build
WORKDIR /build
COPY . /build
RUN mkdir /out && \
    CGO_ENABLED=0 GOARCH=amd64 go build -mod vendor -o /out/adapter pkg/main.go

FROM scratch
COPY --from=build /out/adapter /
ENTRYPOINT ["/adapter", "--logtostderr=true"]
