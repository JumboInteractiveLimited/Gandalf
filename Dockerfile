FROM golang:1.11-stretch
ENV GO111MODULE=on
WORKDIR /go/src/github.com/JumboInteractiveLimited/Gandalf
COPY go.mod ./
RUN echo "Pulling go dependencies" \
	&& go mod download
COPY . .
RUN echo "Testing gandalf" \
	&& go test -v -cover -short github.com/JumboInteractiveLimited/Gandalf/...
