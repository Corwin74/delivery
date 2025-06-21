package courier

import (
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/domain/models/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrInvalidCourierID = errors.New("courier id cannot be nil")
	ErrInvalidName      = errors.New("name cannot be empty")
	ErrInvalidSpeed     = errors.New("speed must be positive")
	ErrOrderNotFound    = errors.New("order not found")
)

const (
	defaultStorageName   = "Cумка"
	defaultStorageVolume = 10
)

type Courier struct {
	id       uuid.UUID
	name     string
	speed    int
	location kernel.Location
	places   []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	if speed <= 0 {
		return nil, ErrInvalidSpeed
	}

	defaultStorage, err := NewStoragePlace(defaultStorageName, defaultStorageVolume)
	if err != nil {
		return nil, err
	}

	return &Courier{
		id:       uuid.New(),
		name:     name,
		speed:    speed,
		location: location,
		places:   []*StoragePlace{defaultStorage},
	}, nil

}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) Places() []*StoragePlace {
	return c.places
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	storagePlace, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.places = append(c.places, storagePlace)
	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) bool {
	if order == nil {
		return false
	}

	volume := order.Volume()

	for _, place := range c.Places() {
		if ok, err := place.CanStore(volume); err == nil && ok {
			return true
		}
	}

	return false
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	for _, place := range c.Places() {
		if ok, err := place.CanStore(order.Volume()); err == nil && ok {
			return place.Store(order.ID(), order.Volume())
		}
	}

	return errors.New("cannot find suitable storage")

}

func (c *Courier) CompleteOrder(order *order.Order) error {
	for _, place := range c.Places() {
		if place.OrderID() != nil && order.ID() == *place.OrderID() {
			return place.Clear(order.ID())
		}
	}

	return ErrOrderNotFound

}

func (c *Courier) CalculateTimeToLocation(location kernel.Location) (float64, error) {
	if location.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("location")
	}

	distance, err := location.DistanceTo(c.Location())
	if err != nil {
		return 0, err
	}
	return float64(distance) / float64(c.Speed()), nil
}

func (c *Courier) Move(target kernel.Location) error {
	if target.IsEmpty() {
		return errs.NewValueIsRequiredError("location")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}

	c.location = newLocation

	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}

	for _, place := range c.Places() {
		if place.OrderID() != nil && *place.OrderID() == orderID {
			return place, nil
		}
	}

	return nil, nil
}
