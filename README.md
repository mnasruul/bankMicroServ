# Microservices Choreography

A demonstration for event sourcing using Go and Kafka example Microservices Choreography.

To run this project:

1. Install Go
2. Install Kafka
3. Run `go mod vendor`
4. Install Redis -> in docker `docker run --name my-first-redis -d redis`
5. Run `go build -o bankMicroServCon.exe && ./bankMicroServCon --act=consumer` to run the program as consumer
6. Run `go build -o bankMicroServPro.exe && ./bankMicroServPro` to run the program as producer

To run testing:

    docker-compose run app ginkgo

Or, you may as well follow the tutorial here: 

Thanks for Adam link:https://github.com/adamnoto/banku
