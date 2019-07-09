FROM golang:1.12 as build-base

WORKDIR /build/app

ADD . /build

RUN CGO_ENABLED=1 go test -mod=readonly -race -coverprofile=.testCoverage.txt -covermode=atomic -coverpkg=$(go list ./... | tr '\n' , | sed 's/,$//') ./...

RUN go tool cover -html=.testCoverage.txt -o coverage.html

RUN go tool cover -func=.testCoverage.txt

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /release/hippo .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o /release/healthcheck "github.com/Soluto/golang-docker-healthcheck/healthcheck"

FROM scratch

ENV PORT 8080

COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-base /release/hippo /hippo
COPY --from=build-base /release/healthcheck /healthcheck

ENTRYPOINT ["/hippo"]

