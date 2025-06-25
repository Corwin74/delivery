package order

import (
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrCannotCompleteNotAssignedOrder   = errors.New("can not complete not assigned order")
	ErrCannotAssignAlreadyAssignedOrder = errors.New("can not assign already assigned order")
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsOutOfRangeError("volume", volume, 1, math.MaxInt)
	}

	return &Order{
		id:       orderID,
		location: location,
		volume:   volume,
		status:   Created,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() int {
	return int(o.status)
}

func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.id == other.id
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if courierID == uuid.Nil {
		return errs.NewValueIsRequiredError("courierID")
	}

	if o.status != Created {
		return ErrCannotAssignAlreadyAssignedOrder
	}

	o.courierID = &courierID
	o.status = Assigned
	return nil
}

func (o *Order) Complete() error {
	if o.status != Assigned {
		return ErrCannotCompleteNotAssignedOrder
	}

	o.status = Completed
	return nil
}
