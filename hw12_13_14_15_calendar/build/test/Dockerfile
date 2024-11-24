FROM golang:1.22 as build

ENV CODE_DIR /go/src/test

WORKDIR ${CODE_DIR}

COPY . .
RUN go mod tidy

#CMD ["go", "test", "-v", "./test"]
CMD ["tail", "-f", "/dev/null"]