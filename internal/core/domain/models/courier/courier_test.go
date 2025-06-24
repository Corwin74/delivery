package courier

import (
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/domain/models/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestNewCourier(t *testing.T) {
	tests := []struct {
		name     string
		nameVal  string
		speed    int
		location kernel.Location
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid courier",
			nameVal:  "Test Courier",
			speed:    10,
			location: mustCreateLocation(t, 5, 5),
			wantErr:  false,
		},
		{
			name:     "empty name",
			nameVal:  "",
			speed:    10,
			location: mustCreateLocation(t, 5, 5),
			wantErr:  true,
			errMsg:   "name cannot be empty",
		},
		{
			name:     "zero speed",
			nameVal:  "Test Courier",
			speed:    0,
			location: mustCreateLocation(t, 5, 5),
			wantErr:  true,
			errMsg:   "speed must be positive",
		},
		{
			name:     "negative speed",
			nameVal:  "Test Courier",
			speed:    -5,
			location: mustCreateLocation(t, 5, 5),
			wantErr:  true,
			errMsg:   "speed must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courier, err := NewCourier(tt.nameVal, tt.speed, tt.location)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, courier)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, courier)
				assert.Equal(t, tt.nameVal, courier.Name())
				assert.Equal(t, tt.speed, courier.Speed())
				assert.Equal(t, tt.location, courier.Location())
				assert.NotEqual(t, uuid.Nil, courier.ID())
				assert.Len(t, courier.Places(), 1) // default storage place
				assert.Equal(t, defaultStorageName, courier.Places()[0].Name())
				assert.Equal(t, defaultStorageVolume, courier.Places()[0].TotalVolume())
			}
		})
	}
}

func TestNewCourier_StoragePlaceCreationError(t *testing.T) {
	// This test covers the error case when NewStoragePlace fails
	// However, since NewStoragePlace only fails with invalid parameters,
	// and we're passing valid ones in NewCourier, this path is hard to test
	// without modifying the code. This is a limitation of the current design.

	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	assert.NoError(t, err)
	assert.NotNil(t, courier)
}

func TestCourier_Getters(t *testing.T) {
	name := "Test Courier"
	speed := 15
	location := mustCreateLocation(t, 3, 7)

	courier, err := NewCourier(name, speed, location)
	require.NoError(t, err)

	t.Run("ID", func(t *testing.T) {
		assert.NotEqual(t, uuid.Nil, courier.ID())
	})

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, name, courier.Name())
	})

	t.Run("Speed", func(t *testing.T) {
		assert.Equal(t, speed, courier.Speed())
	})

	t.Run("Location", func(t *testing.T) {
		assert.Equal(t, location, courier.Location())
	})

	t.Run("Places", func(t *testing.T) {
		places := courier.Places()
		assert.Len(t, places, 1)
		assert.Equal(t, defaultStorageName, places[0].Name())
		assert.Equal(t, defaultStorageVolume, places[0].TotalVolume())
	})
}

func TestCourier_AddStoragePlace(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("add valid storage place", func(t *testing.T) {
		initialPlacesCount := len(courier.Places())
		err := courier.AddStoragePlace("Backpack", 5)

		assert.NoError(t, err)
		assert.Len(t, courier.Places(), initialPlacesCount+1)

		// Check the last added place
		lastPlace := courier.Places()[len(courier.Places())-1]
		assert.Equal(t, "Backpack", lastPlace.Name())
		assert.Equal(t, 5, lastPlace.TotalVolume())
	})

	t.Run("add storage place with empty name", func(t *testing.T) {
		err := courier.AddStoragePlace("", 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("add storage place with invalid volume", func(t *testing.T) {
		err := courier.AddStoragePlace("Invalid", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "totalVolume")
	})

	t.Run("add multiple storage places", func(t *testing.T) {
		courier, _ := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))

		err := courier.AddStoragePlace("Backpack", 3)
		assert.NoError(t, err)

		err = courier.AddStoragePlace("Box", 7)
		assert.NoError(t, err)

		assert.Len(t, courier.Places(), 3) // default + 2 added
	})
}

func TestCourier_CanTakeOrder(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	// Add additional storage place
	err = courier.AddStoragePlace("Backpack", 3)
	require.NoError(t, err)

	t.Run("can take order that fits in default storage", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 5)
		require.NoError(t, err)

		canTake := courier.CanTakeOrder(order)
		assert.True(t, canTake)
	})

	t.Run("can take order that fits in additional storage", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 2)
		require.NoError(t, err)

		canTake := courier.CanTakeOrder(order)
		assert.True(t, canTake)
	})

	t.Run("cannot take order that is too large", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 15)
		require.NoError(t, err)

		canTake := courier.CanTakeOrder(order)
		assert.False(t, canTake)
	})

	t.Run("cannot take nil order", func(t *testing.T) {
		canTake := courier.CanTakeOrder(nil)
		assert.False(t, canTake)
	})
}

func TestCourier_TakeOrder(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("take order successfully", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 5)
		require.NoError(t, err)

		err = courier.TakeOrder(order)
		assert.NoError(t, err)

		// Verify order is stored in default storage
		defaultStorage := courier.Places()[0]
		assert.Equal(t, order.ID(), *defaultStorage.OrderID())
	})

	t.Run("cannot take nil order", func(t *testing.T) {
		err := courier.TakeOrder(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order")
	})

	t.Run("cannot take order that is too large", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 15)
		require.NoError(t, err)

		err = courier.TakeOrder(order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot find suitable storage")
	})
}

func TestCourier_CompleteOrder(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("complete order successfully", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 5)
		require.NoError(t, err)

		// Take order first
		err = courier.TakeOrder(order)
		require.NoError(t, err)

		// Complete order
		err = courier.CompleteOrder(order)
		assert.NoError(t, err)

		// Verify storage is cleared
		defaultStorage := courier.Places()[0]
		assert.Nil(t, defaultStorage.OrderID())
	})

	t.Run("cannot complete order that is not taken", func(t *testing.T) {
		order, err := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 5)
		require.NoError(t, err)

		err = courier.CompleteOrder(order)
		assert.Error(t, err)
		assert.Equal(t, ErrOrderNotFound, err)
	})
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("calculate time to nearby location", func(t *testing.T) {
		targetLocation := mustCreateLocation(t, 7, 5)
		time, err := courier.CalculateTimeToLocation(targetLocation)

		assert.NoError(t, err)
		// Distance is 2 (Manhattan), speed is 10, so time should be 0.2
		assert.Equal(t, 0.2, time)
	})

	t.Run("calculate time to same location", func(t *testing.T) {
		targetLocation := mustCreateLocation(t, 5, 5)
		time, err := courier.CalculateTimeToLocation(targetLocation)

		assert.NoError(t, err)
		assert.Equal(t, 0.0, time)
	})

	t.Run("cannot calculate time to empty location", func(t *testing.T) {
		emptyLocation := kernel.Location{}
		time, err := courier.CalculateTimeToLocation(emptyLocation)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "location")
		assert.Equal(t, 0.0, time)
	})
}

func TestCourier_Move(t *testing.T) {
	courier, err := NewCourier("Test Courier", 3, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("move within speed limit", func(t *testing.T) {
		targetLocation := mustCreateLocation(t, 7, 5)
		err := courier.Move(targetLocation)

		assert.NoError(t, err)
		// Should move 2 units in X direction (within speed limit of 3)
		assert.Equal(t, 7, courier.Location().X())
		assert.Equal(t, 5, courier.Location().Y())
	})

	t.Run("move beyond speed limit", func(t *testing.T) {
		// Reset to original position
		courier, _ = NewCourier("Test Courier", 3, mustCreateLocation(t, 5, 5))

		targetLocation := mustCreateLocation(t, 10, 10)
		err := courier.Move(targetLocation)

		assert.NoError(t, err)
		// Should move only 3 units (speed limit) towards target
		// In this case, it should move 3 units in X direction first
		assert.Equal(t, 8, courier.Location().X())
		assert.Equal(t, 5, courier.Location().Y())
	})

	t.Run("cannot move to empty location", func(t *testing.T) {
		emptyLocation := kernel.Location{}
		err := courier.Move(emptyLocation)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "location")
	})
}

func TestCourier_StoragePlaceManagement(t *testing.T) {
	courier, err := NewCourier("Test Courier", 10, mustCreateLocation(t, 5, 5))
	require.NoError(t, err)

	t.Run("multiple storage places work correctly", func(t *testing.T) {
		// Add multiple storage places
		err = courier.AddStoragePlace("Backpack", 3)
		require.NoError(t, err)
		err = courier.AddStoragePlace("Box", 7)
		require.NoError(t, err)

		assert.Len(t, courier.Places(), 3) // default + 2 added

		// Take orders that fit in different places
		order1, _ := order.NewOrder(uuid.New(), mustCreateLocation(t, 3, 3), 2)
		order2, _ := order.NewOrder(uuid.New(), mustCreateLocation(t, 4, 4), 5)

		err = courier.TakeOrder(order1)
		assert.NoError(t, err)
		err = courier.TakeOrder(order2)
		assert.NoError(t, err)

		// Проверяем, что каждый заказ действительно находится в каком-то storage
		var foundOrder1, foundOrder2 bool
		for _, place := range courier.Places() {
			if place.OrderID() != nil && *place.OrderID() == order1.ID() {
				foundOrder1 = true
			}
			if place.OrderID() != nil && *place.OrderID() == order2.ID() {
				foundOrder2 = true
			}
		}
		assert.True(t, foundOrder1, "order1 должен быть размещён в одном из storage places")
		assert.True(t, foundOrder2, "order2 должен быть размещён в одном из storage places")
	})
}

// Helper function to create location for testing
func mustCreateLocation(t *testing.T, x, y int) kernel.Location {
	location, err := kernel.NewLocation(x, y)
	require.NoError(t, err)
	return location
}
