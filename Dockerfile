FROM golang:1.10-stretch
WORKDIR /go/src/github.com/JumboInteractiveLimited/Gandalf
RUN echo "Pulling go dependencies" \
	&& go get -v -u \
		github.com/fatih/color \
		github.com/JumboInteractiveLimited/jsonpath \
		github.com/jmartin82/mmock \
		github.com/tidwall/gjson \
		gopkg.in/h2non/gock.v1 \
		github.com/eapache/go-resiliency/retrier
COPY . .
RUN echo "Testing gandalf" \
	&& go test -v -cover -short github.com/JumboInteractiveLimited/Gandalf/...
