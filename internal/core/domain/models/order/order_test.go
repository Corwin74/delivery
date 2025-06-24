package order

import (
	"delivery/internal/core/domain/models/kernel"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name     string
		orderID  uuid.UUID
		location kernel.Location
		volume   int
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid order",
			orderID:  uuid.New(),
			location: mustCreateLocation(t, 5, 5),
			volume:   10,
			wantErr:  false,
		},
		{
			name:     "nil order ID",
			orderID:  uuid.Nil,
			location: mustCreateLocation(t, 5, 5),
			volume:   10,
			wantErr:  true,
			errMsg:   "orderID",
		},
		{
			name:     "empty location",
			orderID:  uuid.New(),
			location: kernel.Location{},
			volume:   10,
			wantErr:  true,
			errMsg:   "location",
		},
		{
			name:     "zero volume",
			orderID:  uuid.New(),
			location: mustCreateLocation(t, 5, 5),
			volume:   0,
			wantErr:  true,
			errMsg:   "volume",
		},
		{
			name:     "negative volume",
			orderID:  uuid.New(),
			location: mustCreateLocation(t, 5, 5),
			volume:   -5,
			wantErr:  true,
			errMsg:   "volume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.orderID, tt.location, tt.volume)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, order)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.orderID, order.ID())
				assert.Equal(t, tt.location, order.Location())
				assert.Equal(t, tt.volume, order.Volume())
				assert.Equal(t, Status(Created).String(), order.Status())
				assert.Nil(t, order.CourierID())
			}
		})
	}
}

func TestOrder_Getters(t *testing.T) {
	orderID := uuid.New()
	location := mustCreateLocation(t, 3, 7)
	volume := 15

	order, err := NewOrder(orderID, location, volume)
	require.NoError(t, err)

	t.Run("ID", func(t *testing.T) {
		assert.Equal(t, orderID, order.ID())
	})

	t.Run("Location", func(t *testing.T) {
		assert.Equal(t, location, order.Location())
	})

	t.Run("Volume", func(t *testing.T) {
		assert.Equal(t, volume, order.Volume())
	})

	t.Run("Status", func(t *testing.T) {
		assert.Equal(t, Status(Created).String(), order.Status())
	})

	t.Run("CourierID initially nil", func(t *testing.T) {
		assert.Nil(t, order.CourierID())
	})
}

func TestOrder_Assign(t *testing.T) {
	t.Run("assign courier successfully", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		courierID := uuid.New()
		err = order.Assign(courierID)

		assert.NoError(t, err)
		assert.Equal(t, Status(Assigned).String(), order.Status())
		assert.Equal(t, courierID, *order.CourierID())
	})

	t.Run("cannot assign with nil courier ID", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		err = order.Assign(uuid.Nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "courierID")
		assert.Equal(t, Status(Created).String(), order.Status())
		assert.Nil(t, order.CourierID())
	})

	t.Run("cannot assign already assigned order", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		// First assignment should succeed
		courierID1 := uuid.New()
		err = order.Assign(courierID1)
		require.NoError(t, err)

		// Second assignment should fail
		courierID2 := uuid.New()
		err = order.Assign(courierID2)

		assert.Error(t, err)
		assert.Equal(t, ErrCannotAssignAlreadyAssignedOrder, err)
	})
}

func TestOrder_Complete(t *testing.T) {
	t.Run("complete assigned order successfully", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		// Assign courier first
		courierID := uuid.New()
		err = order.Assign(courierID)
		require.NoError(t, err)

		// Complete order
		err = order.Complete()
		assert.NoError(t, err)
		assert.Equal(t, Status(Completed).String(), order.Status())
	})

	t.Run("cannot complete created order", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		err = order.Complete()
		assert.Error(t, err)
		assert.Equal(t, ErrCannotCompleteNotAssignedOrder, err)
		assert.Equal(t, Status(Created).String(), order.Status())
	})

	t.Run("cannot complete already completed order", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		// Assign and complete order
		courierID := uuid.New()
		err = order.Assign(courierID)
		require.NoError(t, err)
		err = order.Complete()
		require.NoError(t, err)

		// Try to complete again
		err = order.Complete()
		assert.Error(t, err)
		assert.Equal(t, ErrCannotCompleteNotAssignedOrder, err)
		assert.Equal(t, Status(Completed).String(), order.Status())
	})
}

func TestOrder_Equals(t *testing.T) {
	order1, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
	require.NoError(t, err)

	order2, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
	require.NoError(t, err)

	// Create order with same ID as order1
	order3, err := NewOrder(order1.ID(), mustCreateLocation(t, 3, 3), 5)
	require.NoError(t, err)

	t.Run("same order", func(t *testing.T) {
		assert.True(t, order1.Equals(order1))
	})

	t.Run("different orders", func(t *testing.T) {
		assert.False(t, order1.Equals(order2))
	})

	t.Run("same ID different properties", func(t *testing.T) {
		assert.True(t, order1.Equals(order3))
	})

	t.Run("nil order", func(t *testing.T) {
		assert.False(t, order1.Equals(nil))
	})
}

func TestOrder_StatusTransitions(t *testing.T) {
	t.Run("created to assigned to completed", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), mustCreateLocation(t, 5, 5), 10)
		require.NoError(t, err)

		// Initially created
		assert.Equal(t, Status(Created).String(), order.Status())
		assert.Nil(t, order.CourierID())

		// Assign courier
		courierID := uuid.New()
		err = order.Assign(courierID)
		require.NoError(t, err)
		assert.Equal(t, Status(Assigned).String(), order.Status())
		assert.Equal(t, courierID, *order.CourierID())

		// Complete order
		err = order.Complete()
		require.NoError(t, err)
		assert.Equal(t, Status(Completed).String(), order.Status())
		assert.Equal(t, courierID, *order.CourierID()) // CourierID should remain
	})
}

// Helper function to create location for testing
func mustCreateLocation(t *testing.T, x, y int) kernel.Location {
	location, err := kernel.NewLocation(x, y)
	require.NoError(t, err)
	return location
}
