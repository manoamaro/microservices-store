module github.com/manoamaro/microservices-store/products_service

go 1.20

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/manoamaro/microservices-store/commons v0.0.1
	github.com/samber/lo v1.38.1
	go.mongodb.org/mongo-driver v1.12.1
)

replace github.com/manoamaro/microservices-store/commons => ../commons
