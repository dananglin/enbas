# syntax=docker/dockerfile:1
FROM golang:1.23.0

RUN go install github.com/magefile/mage@v1.15.0

ENTRYPOINT ["mage"]
