FROM golang:latest
RUN go get -u github.com/google/go-github/github \
&&	go get -u golang.org/x/oauth2 \
&&  mkdir -p /go/src/github.com/tahirozdemir34/GoMarkdownPusher
COPY . /go/src/github.com/tahirozdemir34/GoMarkdownPusher

WORKDIR /go/src/github.com/tahirozdemir34/GoMarkdownPusher

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/tahirozdemir34/GoMarkdownPusher .
CMD ["./app"]  
