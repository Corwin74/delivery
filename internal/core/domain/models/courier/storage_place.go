package courier

import (
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrCannotStoreOrderInThisStoragePlace = errors.New("cannot store order in this storage place")
	ErrOrderNotStoredInThisPlace          = errors.New("order is not stored in this place")
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if totalVolume <= 0 {
		return nil, errs.NewValueIsOutOfRangeError("totalVolume", totalVolume, 1, math.MaxInt)
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
		orderID:     nil,
	}, nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	if other == nil {
		return false
	}

	return s.id == other.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {
		return false, errs.NewValueIsOutOfRangeError("volume", volume, 1, math.MaxInt)
	}

	if s.isOccupied() {
		return false, nil
	}

	if volume > s.TotalVolume() {
		return false, nil
	}

	return true, nil
}

func (s *StoragePlace) Store(order uuid.UUID, volume int) error {
	if order == uuid.Nil {
		return errs.NewValueIsRequiredError("order")
	}

	if volume <= 0 {
		return errs.NewValueIsOutOfRangeError("volume", volume, 1, math.MaxInt)
	}

	ok, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	if !ok {
		return ErrCannotStoreOrderInThisStoragePlace
	}

	s.orderID = &order
	return nil

}

func (s *StoragePlace) Clear(order uuid.UUID) error {
	if order == uuid.Nil {
		return errs.NewValueIsRequiredError("order")
	}

	if *s.orderID != order {
		return ErrOrderNotStoredInThisPlace
	}

	s.orderID = nil
	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}
