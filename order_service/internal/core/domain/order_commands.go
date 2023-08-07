package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &OrderCreate{} })
}

const (
	OrderCreateCmd = eh.CommandType("order::create")
)

// Static type check that the eventhorizon.Command interface is implemented.
var _ = eh.Command(&OrderCreate{})

type OrderCreate struct {
	ID     uuid.UUID `json:"id"`
	UserId string    `json:"user_id"`
}

func (o *OrderCreate) AggregateID() uuid.UUID          { return o.ID }
func (o *OrderCreate) AggregateType() eh.AggregateType { return OrderAggregateType }
func (o *OrderCreate) CommandType() eh.CommandType     { return OrderCreateCmd }
