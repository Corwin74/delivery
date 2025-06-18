package courier

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStoragePlace(t *testing.T) {
	tests := []struct {
		testName	string	
		name       string
		totalVolume int
		expected   error
	}{
		
		{	
			testName:  "invalid name",
			name:       "",
			totalVolume: 10,
			expected:   errors.New("name cannot be empty"),
		},
		{
			testName:  "invalid total volume",
			name:       "box",
			totalVolume: -1,
			expected:   fmt.Errorf("totalVolume cannot be negative, got: %d", -1),
		},
		{
			testName:  "valid storage place",
			name:       "box",
			totalVolume: 10,
			expected:   nil,
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

func TestCanStore(t *testing.T) {
	storagePlace, err := NewStoragePlace("box", 5)
	assert.Nil(t, err)

	err = storagePlace.Store(uuid.New(), 10)
	assert.ErrorIs(t, err, ErrCannotStoreOrderInThisStoragePlace)

	err = storagePlace.Store(uuid.New(), 5)
	assert.Nil(t, err)

	err = storagePlace.Store(uuid.New(), 1)
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

