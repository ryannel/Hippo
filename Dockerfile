FROM golang:1.12.5 as build-base

WORKDIR /build
	
ADD . /build
RUN if [ $(go mod tidy -v 2>&1 | grep -c unused) != 0 ]; then echo "Unused modules, please run 'go mod tidy'"; exit 1; fi
RUN go fmt ./...

RUN go vet ./...

RUN CGO_ENABLED=1 go test -mod=readonly -race -coverprofile=.testCoverage.txt -covermode=atomic -coverpkg=$(go list ./... | tr '\n' , | sed 's/,$//') ./...

RUN go tool cover -html=.testCoverage.txt -o coverage.html

RUN go tool cover -func=.testCoverage.txt

RUN if test -f ./main.go; then \
  echo "building main.go" \
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /release/testProject main.go; \
fi

RUN if test -d /build/cmd; then \
  for file in $(find /build/cmd -name *.go); \
    do echo "building ${file}"; \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /release/$(basename $(dirname ${file})) ${file}; \
  done \
fi

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o /release/healthcheck "github.com/Soluto/golang-docker-healthcheck/healthcheck"

FROM scratch

ENV PORT 8080

COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-base /release/* /
COPY --from=build-base /release/healthcheck /healthcheck

ENTRYPOINT ["/testProject"]	
