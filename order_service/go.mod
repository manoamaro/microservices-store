module github.com/manoamaro/microservices-store/order_service

go 1.20

require github.com/stretchr/testify v1.8.4
require github.com/looplab/eventhorizon v0.16.0
require github.com/manoamaro/microservices-store/commons v0.0.1

replace github.com/manoamaro/microservices-store/commons => ../commons
