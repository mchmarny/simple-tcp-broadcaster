FROM golang:1.10.1
WORKDIR /go/src/github.com/mchmarny/simple-server/
COPY . .

# restore to pinnned versions of dependancies 
RUN go get github.com/golang/dep/cmd/dep
RUN	dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build . -o simple


FROM scratch
COPY --from=0 /go/src/github.com/mchmarny/simple-server/simple .
ENTRYPOINT ["/simple server --port 3030"]