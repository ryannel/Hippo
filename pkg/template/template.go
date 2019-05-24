package template

import "strings"

func PostgresDeployYaml(POSTGRES_DB string, POSTGRES_USER string, POSTGRES_PASSWORD string) string {
	template := `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: postgresql
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: postgresql
    spec:
      containers:
        - name: postgres
          image: postgres:10.4
          imagePullPolicy: "IfNotPresent"
          resources:
            requests:
              memory: 200Mi
            limits:
              memory: 300Mi
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: {dbName}
            - name: POSTGRES_USER
              value: {user}
            - name: POSTGRES_PASSWORD
              value: {password}
---
apiVersion: v1
kind: Service
metadata:
  name: postgresql
  labels:
    app: postgresql
spec:
  ports:
   - port: 5432
  selector:
   app: postgresql
  type: LoadBalancer
`

	template = strings.Replace(template, "{dbName}", POSTGRES_DB, -1)
	template = strings.Replace(template, "{user}", POSTGRES_USER, -1)
	return strings.Replace(template, "{password}", POSTGRES_PASSWORD, -1)
}

func RabbitDeployYaml(user string, password string) string {
	template := `apiVersion: apps/v1
kind: Deployment
metadata:
  	labels:
    	app: rabbitmq
 	name: rabbitmq
spec:
  	replicas: 1
  	selector:
    	matchLabels:
	      	app: rabbitmq
  	template:
		metadata:
    	  	labels:
	    	    app: rabbitmq
    spec:
      	containers:
        	- image: pliljenberg/rabbitmq
          	name: rabbitmq
          	resources:
	            requests:
    	          	memory: 200Mi
        	    limits:
            	  	memory: 300Mi
          	ports:	
            	- name: management
              		containerPort: 15672
            	- name: rabbit
              		containerPort: 5672
          	env:
            	- name: RABBITMQ_DEFAULT_USER
              		value: {user}
            	- name: RABBITMQ_DEFAULT_PASS
              		value: {password}
---
apiVersion: v1
kind: Service
metadata:
	name: rabbitmq
  	labels:
   		app: rabbitmq
spec:
  	ports:
    	- name: rabbit
      	port: 5672
    	- name: management
      	port: 15672
  	selector:
    	app: rabbitmq
  	type: LoadBalancer`

	template = strings.Replace(template, "{user}", user, -1)
	return strings.Replace(template, "{password}", password, -1)
}

func GenericDeployYaml(projectName string, dockerRegistryUrl string) string {
	template := `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: {projectname}
  name: {projectname}
  annotations:
    kubernetes.io/change-cause: "{timestamp} Deployed commit id: {commit}"
spec:
  replicas: 2
  selector:
    matchLabels:
      app: {projectname}
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {projectname}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: "app"
                  operator: In
                  values:
                  - {projectname}
              topologyKey: kubernetes.io/hostname
      containers:
      - name: {projectname}
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 5
        imagePullPolicy: IfNotPresent
        image: {dockerRegistryUrl}/{projectname}:{commit}
        ports:
        - containerPort: 80
      restartPolicy: Always
---

apiVersion: v1
kind: Service
metadata:
  name: {projectname}
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: {projectname}
  type: ClusterIP`

	template = strings.Replace(template, "{dockerRegistryUrl}", dockerRegistryUrl, -1)
	return strings.Replace(template, "{projectname}", projectName, -1)
}

func GoDockerFile(projectName string) string {
	template := `FROM golang:1.12.5 as build-base

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
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /release/{projectname} main.go; \
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

ENTRYPOINT ["/{projectname}"]	
`

	return strings.Replace(template, "{projectname}", projectName, -1)
}