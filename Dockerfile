FROM golang:1.10.1
WORKDIR /go/src/github.com/mchmarny/simple-tcp-broadcaster
COPY . .

# restore to pinnned versions of dependancies 
RUN go get github.com/golang/dep/cmd/dep
RUN	dep ensure

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o simple 


FROM scratch
COPY --from=0 /go/src/github.com/mchmarny/simple-tcp-broadcaster/simple .
ENTRYPOINT ["/simple server start --port 5050"]