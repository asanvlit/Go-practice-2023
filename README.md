# Golang practice 2023

## Used technologies:

* Golang 1.19
* REST API
* Docker, Docker Compose
* PostgreSQL
* NATS message broker
* Gorilla/mux
* migrate/v4
* Testify

## Project functionality:

* The project includes 3 microservices. 
* The first microservice 'go-auth' has an endpoint for creating a user. 
* The second microservice 'go-users' communicates with the first microservice using NATS and stores in its database the users it receives from the first microservice. 
* The third microservice 'go-scheduler' makes a http request every `n` minutes to get information about new users to the first microservice and also stores the same users in its database.

## Additional features :

* For each microservice, the "Health Pings" mechanism is implemented - if the server stops responding to Pings, then a Gracefull shutdown occurs
* Migrations for databases are used