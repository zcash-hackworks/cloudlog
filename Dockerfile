FROM golang:1.13

ADD . /workspace/cloudlog
WORKDIR /workspace/cloudlog
RUN make build
# Testing
ENTRYPOINT ["/workspace/cloudlog/cloudlog"]
