package domain_test

import (
	"github.com/manoamaro/microservices-store/commons/pkg/event_sourcing"
	"github.com/manoamaro/microservices-store/order_service/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type orderSuite struct {
	suite.Suite
}

func (s *orderSuite) TestGivenCreateEventShouldCreateOrder() {
	data := domain.OrderCreated{
		UserId: "user-id",
	}
	createdEvent := domain.NewOrderEvent(event_sourcing.EmptyAggregateID, 1, data)

	order := domain.Order{}

	err := order.Apply(createdEvent.Event)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), order.UserId, "user-id")
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, &orderSuite{})
}
