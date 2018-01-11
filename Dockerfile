FROM golang:1.9-stretch
WORKDIR /go/src/github.com/JumboInteractiveLimited/Gandalf
RUN echo "Pulling go dependencies" \
	&& go get -v -u \
		github.com/fatih/color \
		github.com/NodePrime/jsonpath \
		github.com/jmartin82/mmock \
		github.com/tidwall/gjson
COPY . .
RUN echo "Testing gandalf" \
	&& go test -v -cover -short github.com/JumboInteractiveLimited/Gandalf/...
