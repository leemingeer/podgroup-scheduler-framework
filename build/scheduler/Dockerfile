# compile
ARG ARCH
FROM golang:1.19

WORKDIR /go/src/sigs.k8s.io/scheduler-plugins
COPY . .
RUN ls -all
ARG ARCH
ARG RELEASE_VERSION
ENV GOPROXY=https://goproxy.cn,direct
RUN RELEASE_VERSION=${RELEASE_VERSION} make build-scheduler.$ARCH

# build
FROM alpine:3.17

COPY --from=0 /go/src/sigs.k8s.io/scheduler-plugins/bin/sample-scheduler /bin/sample-scheduler

WORKDIR /bin
# options will be added in yaml.
CMD ["sample-scheduler"]
