FROM debian:stretch-slim

WORKDIR /

COPY bin/sample-scheduler /usr/local/bin

CMD ["sample-scheduler"]
