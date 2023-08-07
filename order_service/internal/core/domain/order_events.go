package domain

import "github.com/manoamaro/microservices-store/commons/pkg/event_sourcing"

type OrderEvent struct {
	event_sourcing.Event
}

func NewOrderEvent(orderId string, version uint64, data interface{}) OrderEvent {
	return OrderEvent{
		Event: event_sourcing.NewEvent(orderId, "ORDER", version, data, nil),
	}
}
