# requirements
- docker

# How to run
docker-compose up -d 
will bring up 
- order service
- payment service
- database
- messaging broker

The database and messaging queue will take some time to start
until they are ready the order service and payment service will be restarting

# How to use

GET /list will return and array of items availble for purchase
```
curl -XGET -H "Content-type: application/json" 'localhost:8090/list'
```
POST /order will allow you to place an order for an item
```
curl -XPOST -H "Content-type: application/json" -d '{
"customer_id": "1",
"item_id": "1",
"quantity": 2,
"payment_info": "999"
}' 'localhost:8090/order'
```
payment processing is mocked so it only passes for "555"
```
curl -XPOST -H "Content-type: application/json" -d '{
"customer_id": "1",
"item_id": "1",
"quantity": 3,
"payment_info": "555"
}' 'localhost:8090/order'
```
POST /orderstatus will show a specified customers orders past and present
```
curl -XPOST -H "Content-type: application/json" -d '{
"customer_id": "1"
}' 'localhost:8090/orderstatus'
```

# Next steps

create a dockerfile for the database that includes the schema and seed info
this will elimination the messy setup in the main.go file

set up a specific exchange for payments rather than using the default

add a logger for errors
currently all the error go to stderr and are accesible through the docker logs command
but a logger would help by storing the logs to debug issues after the application exits

refactor the database request into the repository model
Right now the handlers do all of the work, after refactoring
the handlers will handle the request
the service layer will handle buisness logic
and the repository layer will be responsible for fetching information from the DB

Use a connection pool for databae connections

