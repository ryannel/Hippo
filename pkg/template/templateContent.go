package template

func GetPostgresDeployYaml() string {
	return `apiVersion: extensions/v1beta1
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
	type: LoadBalancer`
}

func GetRabbitDeployYaml() string {
	return `apiVersion: apps/v1
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
}
