package services

import (
	"delivery/internal/core/domain/models/courier"
	ord "delivery/internal/core/domain/models/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"
)

type OrderDispatcher interface {
	Dispatch(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

type orderDispatcher struct {
}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (d *orderDispatcher) Dispatch(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}

	if len(couriers) == 0 {
		return nil, errs.NewValueIsInvalidError("couriers")
	}

	if order.Status() != ord.Created {
		return nil, errors.New("order is already assigned")
	}

	courier, err := findCourier(order, couriers)
	if err != nil {
		return nil, err
	}

	if err := courier.TakeOrder(order); err != nil {
		return nil, err
	}

	if err := order.Assign(courier.ID()); err != nil {
		return nil, err
	}

	return courier, nil
}

func findCourier(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var bestCourier *courier.Courier
	minTime := math.MaxFloat64

	for _, courier := range couriers {
		if !courier.CanTakeOrder(order) {
			continue
		}

		time, err := courier.CalculateTimeToLocation(order.Location())
		if err != nil {
			return nil, err
		}

		if time < minTime {
			minTime = time
			bestCourier = courier
		}

	}

	if bestCourier == nil {
		return nil, errors.New("courier cannot be found")
	}

	return bestCourier, nil
}
