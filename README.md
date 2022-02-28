# tracker-service
A microservice to track webpage views. Written using gin web framework and uses kafka as message queue and mongodb as document store.

### handlers
http://localhost:6985/awesome-page has a 10X10 pixel which fires/mimicks as pixel fire and gin handlers pushes the view event to kafka topic.
http://localhost:6985/awesome-page/views gives overall stats of the total views and unique users based on the query params passed to the http call.

ex: http://localhost:6985/awesome-page/views?startDate=2022-02-28&endDate=2022-03-01 gives total views and unique users count happened in one day.

and also all the services are loosely coupled using message queues and they can be leveraged to be indeependently deployed as microservice.

server also has json api handler which responds with json paylod of the views 


### pre-requisites
#### 1. Zookeeper
#### 2. Kafka
#### 3. MongoDB
#### 4. go 1.16



### How to run the application
#### brew services start zookeeper
#### brew services start kafka
#### mongod 
#### 1. go build & ./tracker-service (OR)
#### 2. go run main.go graceful_shutdown.go

