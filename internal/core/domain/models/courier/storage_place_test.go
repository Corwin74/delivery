package courier

import (
	"delivery/internal/pkg/errs"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStoragePlace(t *testing.T) {
	tests := []struct {
		testName    string
		name        string
		totalVolume int
		expected    error
	}{

		{
			testName:    "invalid name",
			name:        "",
			totalVolume: 10,
			expected:    errs.NewValueIsRequiredError("name"),
		},
		{
			testName:    "invalid total volume",
			name:        "box",
			totalVolume: -1,
			expected:    errs.NewValueIsOutOfRangeError("totalVolume", -1, 1, math.MaxInt),
		},
		{
			testName:    "valid storage place",
			name:        "box",
			totalVolume: 10,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewStoragePlace(tt.name, tt.totalVolume)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestCanStoreExceedsCapacity(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	ok, err := storagePlace.CanStore(6)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestStore(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	// Try to store order that exceeds capacity
	err = storagePlace.Store(uuid.New(), 10)
	assert.ErrorIs(t, err, ErrCannotStoreOrderInThisStoragePlace)

	// Store order that fits exactly
	orderID1 := uuid.New()
	err = storagePlace.Store(orderID1, 5)
	assert.Nil(t, err)
	assert.Equal(t, orderID1, *storagePlace.OrderID())

	// Try to store another order when already occupied
	err = storagePlace.Store(uuid.New(), 1)
	assert.ErrorIs(t, err, ErrCannotStoreOrderInThisStoragePlace)
}

func TestClearStore(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	orderID := uuid.New()

	err = storagePlace.Store(orderID, 3)
	assert.Nil(t, err)

	err = storagePlace.Clear(uuid.New())
	assert.ErrorIs(t, err, ErrOrderNotStoredInThisPlace)

	err = storagePlace.Clear(orderID)
	assert.Nil(t, err)

	err = storagePlace.Store(uuid.New(), 1)
	assert.Nil(t, err)
}

func TestGetters(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	assert.Equal(t, "box", storagePlace.Name())
	assert.Equal(t, 5, storagePlace.TotalVolume())
	assert.Nil(t, storagePlace.OrderID())

	orderID := uuid.New()
	err = storagePlace.Store(orderID, 3)
	assert.Nil(t, err)

	assert.Equal(t, orderID, *storagePlace.OrderID())
}

func TestIsOccupied(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	assert.False(t, storagePlace.isOccupied())

	orderID := uuid.New()
	err = storagePlace.Store(orderID, 3)
	assert.Nil(t, err)

	assert.True(t, storagePlace.isOccupied())

	err = storagePlace.Clear(orderID)
	assert.Nil(t, err)

	assert.False(t, storagePlace.isOccupied())
}
