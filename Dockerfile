FROM golang:1.13

ADD . /workspace/cloudlog
WORKDIR /workspace/cloudlog
RUN make build
ENTRYPOINT ["/workspace/cloudlog/cloudlog"]
