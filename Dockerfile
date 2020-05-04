#build stage
FROM golang:latest AS builder
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM cockroachdb/cockroach:latest
LABEL Name=cockroachdb-backup Version=0.0.2 maintainer="Yash Thakkar<thakkaryash94@gmail.com>"
RUN mkdir /data /cockroach-certs
COPY --from=builder /go/bin/app /
VOLUME [ "/data" ]
VOLUME [ "/cockroach-certs" ]
EXPOSE 9000
ENTRYPOINT ["/bin/bash"]
CMD [ "-c", "/app" ]
