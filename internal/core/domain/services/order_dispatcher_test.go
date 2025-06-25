package services

import (
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
	ord "delivery/internal/core/domain/models/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewOrderDispatcher()

	t.Run("successfully dispatch order to nearest courier", func(t *testing.T) {
		// Создаем заказ
		order, err := ord.NewOrder(uuid.New(), mustCreateLocation(t, 10, 10), 5)
		require.NoError(t, err)

		// Создаем курьеров на разных расстояниях
		courier1, err := courier.NewCourier("Courier 1", 10, mustCreateLocation(t, 5, 5)) // ближе
		require.NoError(t, err)

		courier2, err := courier.NewCourier("Courier 2", 10, mustCreateLocation(t, 3, 3)) // дальше
		require.NoError(t, err)

		couriers := []*courier.Courier{courier1, courier2}

		// Диспетчеризируем заказ
		assignedCourier, err := dispatcher.Dispatch(order, couriers)

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, assignedCourier)
		assert.Equal(t, courier1.ID(), assignedCourier.ID()) // должен быть выбран ближайший
		assert.Equal(t, ord.Assigned, order.Status())
		assert.Equal(t, courier1.ID(), *order.CourierID())
	})

	t.Run("return error when order is nil", func(t *testing.T) {
		courier1, err := courier.NewCourier("Courier 1", 10, mustCreateLocation(t, 5, 5))
		require.NoError(t, err)

		couriers := []*courier.Courier{courier1}

		assignedCourier, err := dispatcher.Dispatch(nil, couriers)

		assert.Error(t, err)
		assert.Nil(t, assignedCourier)
		assert.Contains(t, err.Error(), "order")
	})

	t.Run("return error when couriers list is empty", func(t *testing.T) {
		order, err := ord.NewOrder(uuid.New(), mustCreateLocation(t, 10, 10), 5)
		require.NoError(t, err)

		couriers := []*courier.Courier{}

		assignedCourier, err := dispatcher.Dispatch(order, couriers)

		assert.Error(t, err)
		assert.Nil(t, assignedCourier)
		assert.Contains(t, err.Error(), "couriers")
	})

	t.Run("return error when order is already assigned", func(t *testing.T) {
		order, err := ord.NewOrder(uuid.New(), mustCreateLocation(t, 10, 10), 5)
		require.NoError(t, err)

		// Назначаем заказ курьеру
		courierID := uuid.New()
		err = order.Assign(courierID)
		require.NoError(t, err)

		courier1, err := courier.NewCourier("Courier 1", 10, mustCreateLocation(t, 5, 5))
		require.NoError(t, err)
		couriers := []*courier.Courier{courier1}

		assignedCourier, err := dispatcher.Dispatch(order, couriers)

		assert.Error(t, err)
		assert.Nil(t, assignedCourier)
		assert.Contains(t, err.Error(), "already assigned")
	})

	t.Run("return error when no courier can take order", func(t *testing.T) {
		// Создаем заказ с большим объемом
		order, err := ord.NewOrder(uuid.New(), mustCreateLocation(t, 10, 10), 20)
		require.NoError(t, err)

		// Создаем курьера с недостаточным местом (по умолчанию 10)
		courier1, err := courier.NewCourier("Courier 1", 10, mustCreateLocation(t, 5, 5))
		require.NoError(t, err)

		couriers := []*courier.Courier{courier1}

		assignedCourier, err := dispatcher.Dispatch(order, couriers)

		assert.Error(t, err)
		assert.Nil(t, assignedCourier)
		assert.Contains(t, err.Error(), "cannot be found")
	})
}

func mustCreateLocation(t *testing.T, x, y int) kernel.Location {
	location, err := kernel.NewLocation(x, y)
	require.NoError(t, err)
	return location
}
