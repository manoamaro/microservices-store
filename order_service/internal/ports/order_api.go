package ports

import (
	"context"
)

type OrderApi interface {
	GetOrderHandler(c context.Context)
}
