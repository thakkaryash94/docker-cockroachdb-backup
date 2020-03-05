
#build stage
FROM golang:latest AS builder
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM cockroachdb/cockroach:v19.2.4
LABEL Name=cockroachdb-backup Version=0.0.1 maintainer="Yash Thakkar<thakkaryash94@gmail.com>"
RUN mkdir /data
COPY --from=builder /go/bin/app /
VOLUME [ "/data" ]
ENTRYPOINT ["/bin/bash"]
CMD [ "-c", "/app" ]
